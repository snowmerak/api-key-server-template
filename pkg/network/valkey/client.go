package valkey

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

type Option struct {
	address  []string
	username string
	password string
	db       int
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

func (o *Option) WithDB(db int) *Option {
	o.db = db
	return o
}

type Client struct {
	option *Option
	client valkey.Client
}

func NewClient(ctx context.Context, option *Option) (*Client, error) {
	opt := valkey.ClientOption{
		InitAddress: option.address,
	}
	if option.username != "" {
		opt.Username = option.username
	}
	if option.password != "" {
		opt.Password = option.password
	}
	if option.db != 0 {
		opt.SelectDB = option.db
	}

	cli, err := valkey.NewClient(opt)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	context.AfterFunc(ctx, func() {
		cli.Close()
	})

	return &Client{
		option: option,
		client: cli,
	}, nil
}
