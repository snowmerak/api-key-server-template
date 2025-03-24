package httpx

import (
	"context"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultClientSideRequestTimeout = 5 * time.Second
)

type ClientOption struct {
	timeout time.Duration
	baseUrl string
}

func NewClientOption() *ClientOption {
	return &ClientOption{
		timeout: DefaultClientSideRequestTimeout,
		baseUrl: "http://localhost:3030",
	}
}

func (o *ClientOption) WithTimeout(timeout time.Duration) *ClientOption {
	o.timeout = timeout
	return o
}

func (o *ClientOption) WithBaseUrl(baseUrl string) *ClientOption {
	o.baseUrl = strings.TrimSuffix(baseUrl, "/")
	return o
}

type Client struct {
	option *ClientOption
	client *http.Client
}

func NewClient(ctx context.Context, option ClientOption) *Client {
	client := &http.Client{
		Timeout: option.timeout,
	}

	context.AfterFunc(ctx, func() {
		client.CloseIdleConnections()
	})

	return &Client{
		option: &option,
		client: client,
	}
}
