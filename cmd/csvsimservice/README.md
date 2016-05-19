Copyright (c) 2016 Johann Höchtl
See LICENSE for license.

This is a web service which compares the similarity of CSV headers. It accepts

on route `/compare` via  PUT:

    type CSVPairCompareRequest struct {
    	Baseline *string
    	Compare  []string
    }

a JSON-encoded body (application/json) with
* a string named `Baseline` which contains a CSV header;
* an array of strings `Compare` of CSV headers against which `Baseline` will be compared.

In case of success (HTTP status code 200), the web service produces response of type  application/json

    type CSVPairCompareResponse struct {
    	CSVPairCompareRequest
    	Response struct {
    		CompareResult []int
    	}
    }


with
* `CSVPairCompareRequest` the input to the web serive copied to the response;
* `CompareResult` an array of integers as the result of measuring the similarity of CSV headers provided in `CSVPairCompareRequest`.

Additionally, these query parameters are accepted:
* `comma`: character which is used to separates CSV fields, defaults to `,`
* `comment`: character to use for comments in the CSV file preceding the header; defaults to none.

This service will read from environment variable `PORT` the port on which it listens for incoming requests and defaults to port 5000.


This service can be either consumed directly after installing [Golang](https://golang.org/) and running

    go get github.com/the42/setsim/cmd/csvsimservice

    # start service
    ./csvsimservice
    # start on another port but the default
    PORT=5001 ./csvsimservice

or using [Docker](https://www.docker.com/)

    docker pull the42/csvsimservice
    docker run -it -p 5000:5000 --rm --name mycsvcomp csvsimservice

which will map exposed port 5000 and make it availavble to the host again as port 5000. Refer to the docker documentation if you want to use [another port](https://docs.docker.com/engine/reference/run/#expose-incoming-ports).

In case of bad input (missing parameters or wrong paramters) the service will respond with http status code 400.  
In case of any other error (like malformed request or internal server error) the service will respond with http status code 500. In both cases the response will be of type `text/plain`.

An example how to use this service is provided in file `in.json` which will be used by the script `test.cmd` which uses the command line tool `curl` to issue a PUT-request to the service.
