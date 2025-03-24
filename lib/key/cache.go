package key

import (
	"context"
	"time"
)

type CacheData StoreData

type CacheOption struct {
	ExpireAt *time.Time
}

type Cache interface {
	Load(ctx context.Context, namespace string, key string) (*CacheData, error)
	Save(ctx context.Context, namespace string, key string, data *CacheData, option *CacheOption) error
	Delete(ctx context.Context, namespace string, key string) error
}
