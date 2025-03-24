package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/snowmerak/api-key-server-template/lib/key"
	"github.com/snowmerak/api-key-server-template/pkg/client/postgres/queries"
)

const (
	DefaultApiKeyExpire = 30 * 24 * 60 * 60 * time.Second
)

func (db *DB) Load(ctx context.Context, namespace string, apiKey string) (*key.StoreData, error) {
	resp, err := db.queries.GetApiKey(ctx, queries.GetApiKeyParams{
		Namespace: namespace,
		ApiKey:    apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting key from API: %v", err)
	}

	sd := &key.StoreData{
		Owner:      resp.Owner,
		Service:    resp.Service,
		Permission: resp.Permissions,
		Payload:    resp.Payload,
	}

	return sd, nil
}

func (db *DB) Save(ctx context.Context, namespace string, apiKey string, data *key.StoreData, option *key.SaveOption) error {
	expiresAt := time.Now().Add(DefaultApiKeyExpire)
	if option.ExpireAt != nil {
		expiresAt = *option.ExpireAt
	}
	option.ExpireAt = &expiresAt

	params := queries.CreateApiKeyParams{
		Namespace:   namespace,
		ApiKey:      apiKey,
		Owner:       data.Owner,
		Service:     data.Service,
		Permissions: data.Permission,
		Payload:     data.Payload,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: true,
		},
	}
	if _, err := db.queries.CreateApiKey(ctx, params); err != nil {
		return fmt.Errorf("error creating key: %v", err)
	}

	return nil
}

func (db *DB) Delete(ctx context.Context, namespace string, apiKey string) error {
	if err := db.queries.DeleteApiKey(ctx, queries.DeleteApiKeyParams{
		Namespace: namespace,
		ApiKey:    apiKey,
	}); err != nil {
		return fmt.Errorf("error deleting key: %v", err)
	}

	return nil
}
