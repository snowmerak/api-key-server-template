package valkey

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/snowmerak/api-key-server-template/lib/key"
)

var _ key.Cache = (*Client)(nil)

const (
	defaultPrefix = "api-key"
)

const (
	defaultClientSideCacheTTL = 5 * time.Minute
)

var (
	ErrorKeyDataIsNil = fmt.Errorf("key data is nil")
)

func makeKey(namespace string, key string) string {
	return defaultPrefix + namespace + ":" + key
}

func (c *Client) Load(ctx context.Context, namespace string, apiKey string) (*key.CacheData, error) {
	k := makeKey(namespace, apiKey)
	resp, err := c.client.DoCache(ctx, c.client.B().Get().Key(k).Cache(), defaultClientSideCacheTTL).AsBytes()
	if err != nil {
		return nil, fmt.Errorf("error getting key from cache: %w", err)
	}
	if resp == nil {
		return nil, ErrorKeyDataIsNil
	}

	cd := &key.CacheData{}
	if err := json.Unmarshal(resp, cd); err != nil {
		return nil, fmt.Errorf("error unmarshalling cache data: %w", err)
	}

	return cd, nil
}

func (c *Client) Save(ctx context.Context, namespace string, apiKey string, cacheData *key.CacheData, option *key.CacheOption) error {
	data, err := json.Marshal(cacheData)
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	k := makeKey(namespace, apiKey)
	if err := c.client.Do(ctx, c.client.B().Set().Key(k).Value(valkey.BinaryString(data)).Build()).Error(); err != nil {
		return fmt.Errorf("error setting key: %w", err)
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, namespace string, apiKey string) error {
	if err := c.client.Do(ctx, c.client.B().Del().Key(makeKey(namespace, apiKey)).Build()).Error(); err != nil {
		return fmt.Errorf("error deleting key: %w", err)
	}

	return nil
}
