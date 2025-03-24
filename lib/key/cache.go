package key

import "time"

type CacheData StoreData

type CacheOption struct {
	ExpireAt *time.Time
}

type Cache interface {
	Load(namespace string, key string) (*CacheData, error)
	Save(namespace string, key string, data *CacheData, option *CacheOption) error
	Delete(namespace string, key string) error
}
