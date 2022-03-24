package margin

import (
	"net/http"

	"github.com/tocoteron/gmo-coin-go/api/v1/http/endpoint"
	"github.com/tocoteron/gmo-coin-go/api/v1/http/endpoint/common"
)

type Margin struct {
	common.CommonResponse
	Data MarginData `json:"data"`
}

type MarginData struct {
	ActualProfitLoss string `json:"actualProfitLoss"`
	AvailableAmount  string `json:"availableAmount"`
	Margin           string `json:"margin"`
	MarginCallStatus string `json:"marginCallStatus"`
	MarginRatio      string `json:"marginRatio"`
	ProfitLoss       string `json:"profitLoss"`
}

var Endpoint = endpoint.NewEndpoint[any, any, Margin](http.MethodGet, "/v1/account/margin")
