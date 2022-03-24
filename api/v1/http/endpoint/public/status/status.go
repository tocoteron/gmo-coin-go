package status

import (
	"net/http"

	"github.com/tocoteron/gmo-coin-go/api/v1/http/endpoint"
	"github.com/tocoteron/gmo-coin-go/api/v1/http/endpoint/common"
)

type Status struct {
	common.CommonResponse
	Data StatusData `json:"data"`
}

type StatusData struct {
	Status string `json:"status"`
}

var Endpoint = endpoint.NewEndpoint[any, any, Status](http.MethodGet, "/v1/status")
