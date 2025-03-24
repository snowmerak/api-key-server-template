package nats

import (
	"context"
	"sync/atomic"

	"github.com/nats-io/nats.go"
)

type Option struct {
	address  []string
	username string
	password string
}

func NewOption(address ...string) *Option {
	return &Option{
		address: address,
	}
}

func (o *Option) WithUsername(username string) *Option {
	o.username = username
	return o
}

func (o *Option) WithPassword(password string) *Option {
	o.password = password
	return o
}

type Client struct {
	option       *Option
	conn         *nats.Conn
	subscription atomic.Pointer[nats.Subscription]
}

func NewClient(ctx context.Context, option *Option) (*Client, error) {
	opt := make([]nats.Option, 0)
	if option.username != "" {
		opt = append(opt, nats.UserInfo(option.username, option.password))
	}

	conn, err := nats.Connect(option.address[0], opt...)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		conn.Close()
	})

	return &Client{
		option: option,
		conn:   conn,
	}, nil
}
