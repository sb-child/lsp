package mtools

import (
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
)

func NewMyHttpClient() *httpclient.Client {
	return httpclient.NewClient(httpclient.WithHTTPTimeout(2 * time.Second))
}

func init() {
}
