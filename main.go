package main

import (
	"flag"
	"fmt"

	v1 "github.com/tocoteron/gmo-coin-go/api/v1"
	"github.com/tocoteron/gmo-coin-go/api/v1/auth"
)

func main() {
	API_KEY := flag.String("key", "", "API Key")
	API_SECRET := flag.String("secret", "", "API Secret")
	flag.Parse()

	c := v1.NewHTTPClient(&v1.ClientOpts{
		AuthConfig: &auth.AuthConfig{
			APIKey:    *API_KEY,
			APISecret: *API_SECRET,
		},
	})

	status, httpResp, err := c.Status()
	if err != nil {
		fmt.Println(httpResp, err)
	}
	fmt.Println(status)

	margin, httpResp, err := c.Margin()
	if err != nil {
		fmt.Println(httpResp, err)
	}
	fmt.Println(margin)

}
