package httpx

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"google.golang.org/protobuf/proto"

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

func (s *ApiServer) ListenAndServe(ctx context.Context, addr string, tlsConfig *tls.Config) error {
	s.Server.handler.HandleFunc(DefaultAuthAPIRoute, func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		defer request.Body.Close()

		data, err := io.ReadAll(request.Body)
		if err != nil {
			log.Printf("failed to read body: %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		req := &v1.AuthorizerRequest{}
		if err := proto.Unmarshal(data, req); err != nil {
			log.Printf("failed to unmarshal: %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := s.Reply(ctx, req)
		data, err = proto.Marshal(resp)
		if err != nil {
			log.Printf("failed to marshal: %v", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Write(data)
	})

	return s.ListenAndServe(ctx, addr, tlsConfig)
}
