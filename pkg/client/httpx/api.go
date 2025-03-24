package httpx

import (
	"context"
	"log"
	"net/http"

	v1 "github.com/snowmerak/api-key-server-template/gen/api/authorizer/v1"
	"github.com/snowmerak/api-key-server-template/lib/api"
	"github.com/snowmerak/api-key-server-template/lib/key"
)

const (
	DefaultAuthAPIRoute = "/api/key/verification"
)

var _ api.SyncServer = (*ApiServer)(nil)

type ApiServer struct {
	Server
	store key.ReadOnlyStore
	cache key.Cache
}

func NewApiServer() *ApiServer {
	return &ApiServer{
		Server: Server{
			server: &http.Server{},
		},
	}
}

func (s *ApiServer) Reply(ctx context.Context, request *v1.AuthorizerRequest) *v1.AuthorizerResponse {
	cd, err := s.cache.Load(ctx, request.GetNamespace(), request.GetToken())
	if err == nil {
		return &v1.AuthorizerResponse{
			Response: &v1.AuthorizerResponse_AuthorizedData{
				AuthorizedData: &v1.AuthorizedData{
					Owner:      cd.Owner,
					Service:    cd.Service,
					Permission: cd.Permission,
					Payload:    cd.Payload,
				},
			},
		}
	}

	sd, err := s.store.Load(ctx, request.GetNamespace(), request.GetToken())
	if err == nil {
		return &v1.AuthorizerResponse{
			Response: &v1.AuthorizerResponse_AuthorizedData{
				AuthorizedData: &v1.AuthorizedData{
					Owner:      sd.Owner,
					Service:    sd.Service,
					Permission: sd.Permission,
					Payload:    sd.Payload,
				},
			},
		}
	}

	log.Printf("failed to get data: %v", err)

	return &v1.AuthorizerResponse{
		Response: &v1.AuthorizerResponse_Error{
			Error: &v1.Error{
				Code:    400,
				Message: "failed to get api key data",
			},
		},
	}
}
