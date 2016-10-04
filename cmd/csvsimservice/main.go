package main

// Copyright (c) Johann Höchtl 2016
//
// See LICENSE for License

// RESTful service to check for the similarity of CSV headers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/the42/setsim"
)

// CSVPairCompareRequest document
type CSVPairCompareRequest struct {
	Baseline *string  `description:"The CSV header the other headers should be compared to"`
	Compare  []string `description:"The CSV headers which are all compared against Baseline"`
}

// CSVPairCompareResponse document
type CSVPairCompareResponse struct {
	CSVPairCompareRequest `description:"verbatim copy of request struct"`
	Response              struct {
		CompareResult []int `description:"Contains for every CSV header in CSVPairCompareRequest.Compare the comparison result"`
	}
}

func compareCSVHeader(request *restful.Request, response *restful.Response) {
	csvset := CSVPairCompareRequest{}
	if err := request.ReadEntity(&csvset); err != nil {
		logresponse(response, http.StatusBadRequest, fmt.Sprintf("unable to parse request: %s", err.Error()))
		return
	}

	if csvset.Baseline == nil || csvset.Compare == nil {
		logresponse(response, http.StatusBadRequest, "insufficient data to compare")
		return
	}

	basereader := csv.NewReader(strings.NewReader(*csvset.Baseline))

	if comma := request.QueryParameter("comma"); len(comma) != 0 {
		basereader.Comma = rune(comma[0])
	}
	if comment := request.QueryParameter("comment"); len(comment) != 0 {
		basereader.Comment = rune(comment[0])
	}

	baserecords, err := basereader.ReadAll()
	if err != nil {
		logresponse(response, http.StatusInternalServerError, err.Error())
		return
	}
	if len(baserecords) == 0 {
		logresponse(response, http.StatusInternalServerError, "no baseline header")
		return
	}

	result := CSVPairCompareResponse{CSVPairCompareRequest: csvset}
	for _, val := range csvset.Compare {
		comparereader := csv.NewReader(strings.NewReader(val))
		comparereader.Comma = basereader.Comma
		comparereader.Comment = basereader.Comment
		comparerecords, err := comparereader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		result.Response.CompareResult = append(result.Response.CompareResult, setsim.StringDistance(baserecords[0], comparerecords[0]))
	}
	response.WriteAsJson(result)
}

func logresponse(resp *restful.Response, code int, message string) {
	resp.WriteErrorString(code, message)
	log.Print(message)
}

func main() {
	ws := new(restful.WebService).
		Produces(restful.MIME_JSON).
		Consumes(restful.MIME_JSON)

	//BEGIN: CORS support
	/*
		cors := restful.CrossOriginResourceSharing{
			ExposeHeaders:  []string{"X-My-Header"},
			AllowedHeaders: []string{"Content-Type", "Accept"},
			AllowedMethods: []string{"GET", "POST", "PUT"},
			CookiesAllowed: false,
			Container:      restful.DefaultContainer}

		restful.DefaultContainer.Filter(cors.Filter)
		// Add container filter to respond to OPTIONS
		restful.DefaultContainer.Filter(restful.DefaultContainer.OPTIONSFilter)
	*/
	//END: CORS support

	ws.Route(ws.PUT("/compare").
		To(compareCSVHeader).
		Produces(restful.MIME_JSON).
		Consumes(restful.MIME_JSON).
		Param(ws.QueryParameter("comma", "CSV field separator character").
			DefaultValue(",").
			Required(false).
			DataType("string")).
		Param(ws.QueryParameter("comment", "CSV comment character").
			DefaultValue("").
			Required(false).
			DataType("string")).
		Reads(CSVPairCompareRequest{}).
		Returns(http.StatusOK, "", CSVPairCompareResponse{}).
		Returns(http.StatusInternalServerError, "", nil).
		Returns(http.StatusBadRequest, "", nil))
	restful.Add(ws)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	config := swagger.Config{
		WebServices:     restful.DefaultContainer.RegisteredWebServices(),
		WebServicesUrl:  "http://localhost",
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "./swagger-ui/dist"}
	swagger.RegisterSwaggerService(config, restful.DefaultContainer)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
