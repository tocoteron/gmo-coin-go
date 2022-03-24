package common

import "time"

type CommonResponse struct {
	Status       int       `json:"status"`
	ResponseTime time.Time `json:"responsetime"`
}
