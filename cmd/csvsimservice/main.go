package main // Copyright (c) Johann Höchtl 2016
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
	"github.com/the42/setsim"
)

// CSVPairCompareRequest document
type CSVPairCompareRequest struct {
	Baseline *string
	Compare  []string
}

// CSVPairCompareResponse document
type CSVPairCompareResponse struct {
	CSVPairCompareRequest
	Response struct {
		CompareResult []int
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

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
