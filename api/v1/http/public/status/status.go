package status

import (
	"net/http"

	"github.com/tocoteron/gmo-coin-go/api/v1/http/common"
)

type Status struct {
	common.CommonResponse
	Data StatusData `json:"data"`
}

type StatusData struct {
	Status string `json:"status"`
}

const (
	Method = http.MethodGet
	Path   = "/v1/status"
)
