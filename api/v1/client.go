package v1

import (
	"github.com/tocoteron/gmo-coin-go/api/v1/auth"
	http "github.com/tocoteron/gmo-coin-go/api/v1/http/client"
)

type ClientOpts struct {
	AuthConfig *auth.AuthConfig
}

func NewHTTPClient(opts *ClientOpts) *http.Client {
	return http.NewClient(&http.ClientOpts{AuthConfig: opts.AuthConfig})
}
