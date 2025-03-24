-- name: GetApiKey :one
SELECT * FROM apiKeys WHERE namespace = $1 AND api_key = $2 AND expired = FALSE;

-- name: GetApiKeysByOwner :many
SELECT * FROM apiKeys WHERE namespace = $1 AND owner = $2 AND expired = FALSE;

-- name: GetApiKeysByService :many
SELECT * FROM apiKeys WHERE namespace = $1 AND service = $2 AND expired = FALSE;

-- name: CreateApiKey :one
INSERT INTO apikeys (namespace, api_key, owner, service, permissions, payload, expired, expires_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: UpdateApiKey :one
UPDATE apikeys SET owner = $2, service = $3, permissions = $4, payload = $5, expired = $6, expires_at = $7, updated_at = now() WHERE api_key = $1 RETURNING *;

-- name: ExpireApiKey :one
UPDATE apikeys SET expired = TRUE AND updated_at = now() WHERE api_key = $1 RETURNING *;

-- name: ExpireApiKeysByOwner :many
UPDATE apikeys SET expired = TRUE AND updated_at = now() WHERE owner = $1 RETURNING *;

-- name: ExpireApiKeysByService :many
UPDATE apikeys SET expired = TRUE AND updated_at = now() WHERE service = $1 RETURNING *;

-- name: ExpireExpiredApiKeys :many
UPDATE apikeys SET expired = TRUE AND updated_at = now() WHERE expires_at < $1 RETURNING *;

-- name: DeleteApiKey :exec
DELETE FROM apikeys WHERE namespace = $1 AND api_key = $2;

-- name: DeleteApiKeysByOwner :exec
DELETE FROM apikeys WHERE namespace = $1 AND owner = $2;

-- name: DeleteApiKeysByService :exec
DELETE FROM apikeys WHERE namespace = $1 AND service = $2;

-- name: DeleteExpiredApiKeys :exec
DELETE FROM apikeys WHERE namespace = $1 AND expired = TRUE;
