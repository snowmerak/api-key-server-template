package key

import (
	"context"
	"time"
)

type StoreData struct {
	Owner      string `json:"owner,omitempty"`
	Service    string `json:"service,omitempty"`
	Permission string `json:"permission,omitempty"`
	Payload    []byte `json:"payload,omitempty"`
}

type SaveOption struct {
	ExpireAt *time.Time
}

type Store interface {
	Load(ctx context.Context, namespace string, key string) (*StoreData, error)
	Save(ctx context.Context, namespace string, key string, data *StoreData, option *SaveOption) error
	Delete(ctx context.Context, namespace string, key string) error
}
