package nats

import (
	"context"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	v1 "github.com/snowmerak/api-key-server-template/gen/api/authorizer/v1"
)

const (
	ApiSubject = "api.authorizer.v1"
)

func (c *Client) Subscribe(ctx context.Context, handler func(request *v1.AuthorizerRequest) *v1.AuthorizerResponse) error {
	if c.subscription.Load() != nil {
		return nil
	}

	sub, err := c.conn.Subscribe(ApiSubject, func(msg *nats.Msg) {
		req := &v1.AuthorizerRequest{}
		if err := proto.Unmarshal(msg.Data, req); err != nil {
			log.Printf("failed to unmarshal authorizer request: %v", err)
			return
		}

		resp := handler(req)
		if resp == nil {
			log.Printf("failed to handle authorizer request: %v", msg)
			return
		}

		data, err := proto.Marshal(resp)
		if err != nil {
			log.Printf("failed to marshal authorizer response: %v", err)
			return
		}

		if err := msg.Respond(data); err != nil {
			log.Printf("failed to respond authorizer response: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %v", ApiSubject, err)
	}

	if !c.subscription.CompareAndSwap(nil, sub) {
		sub.Drain()
		sub.Unsubscribe()
	}

	context.AfterFunc(ctx, func() {
		c.subscription.Store(nil)
		sub.Drain()
		sub.Unsubscribe()
	})

	return nil
}

func (c *Client) Request(ctx context.Context, request *v1.AuthorizerRequest) (*v1.AuthorizerResponse, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal authorizer request: %v", err)
	}

	msg, err := c.conn.RequestWithContext(ctx, ApiSubject, data)
	if err != nil {
		return nil, fmt.Errorf("failed to request authorizer response: %v", err)
	}

	resp := &v1.AuthorizerResponse{}
	if err := proto.Unmarshal(msg.Data, resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal authorizer response: %v", err)
	}

	return resp, nil
}
