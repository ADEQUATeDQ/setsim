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
	Comma    *string
	Comment  *string
	Baseline *string
	Compare  *string
}

// CSVPairCompareResponse document
type CSVPairCompareResponse struct {
	CSVPairCompareRequest
	CompareResult []int `json:"compareResult"`
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
	comparereader := csv.NewReader(strings.NewReader(*csvset.Compare))

	if comma := csvset.Comma; comma != nil {
		basereader.Comma = rune((*comma)[0])
		comparereader.Comma = basereader.Comma
	}
	if comment := csvset.Comment; comment != nil {
		basereader.Comment = rune((*comment)[0])
		comparereader.Comment = basereader.Comment
	}

	baserecords, err := basereader.ReadAll()
	if err != nil {
		logresponse(response, http.StatusInternalServerError, err.Error())
		return
	}
	if len(baserecords) == 0 {
		logresponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	comparerecords, err := comparereader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	if len(comparerecords) == 0 {
		logresponse(response, http.StatusInternalServerError, "no baseline header")
		return
	}

	result := CSVPairCompareResponse{CSVPairCompareRequest: csvset}
	for _, val := range comparerecords {
		result.CompareResult = append(result.CompareResult, setsim.StringDistance(baserecords[0], val))
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

	ws.Route(ws.PUT("/compare").To(compareCSVHeader))
	restful.Add(ws)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
