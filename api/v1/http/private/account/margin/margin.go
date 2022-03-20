package margin

import (
	"net/http"

	"github.com/tocoteron/gmo-coin-go/api/v1/http/common"
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

const (
	Method = http.MethodGet
	Path   = "/v1/account/margin"
)
