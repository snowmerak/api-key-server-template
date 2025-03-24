package api

import (
	"context"

	v1 "github.com/snowmerak/api-key-server-template/gen/api/authorizer/v1"
)

type AsyncServer interface {
	Subscribe(ctx context.Context, handler func(request *v1.AuthorizerRequest) *v1.AuthorizerResponse) error
}

type SyncServer interface {
	Reply(ctx context.Context, request *v1.AuthorizerRequest) *v1.AuthorizerResponse
}

type Client interface {
	Request(ctx context.Context, request *v1.AuthorizerRequest) (*v1.AuthorizerResponse, error)
}
