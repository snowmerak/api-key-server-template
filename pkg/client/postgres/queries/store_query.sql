-- name: GetApiKey :one
SELECT * FROM apiKeys WHERE api_key = $1 AND expired = FALSE;

-- name: GetApiKeysByOwner :many
SELECT * FROM apiKeys WHERE owner = $1 AND expired = FALSE;

-- name: GetApiKeysByService :many
SELECT * FROM apiKeys WHERE service = $1 AND expired = FALSE;

-- name: CreateApiKey :one
INSERT INTO apikeys (api_key, owner, service, permissions, payload, expired, expires_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;

-- name: UpdateApiKey :one
UPDATE apikeys SET owner = $2, service = $3, permissions = $4, payload = $5, expired = $6, expires_at = $7, updated_at = $8 WHERE api_key = $1 RETURNING *;

-- name: ExpireApiKey :one
UPDATE apikeys SET expired = TRUE WHERE api_key = $1 RETURNING *;

-- name: ExpireApiKeysByOwner :many
UPDATE apikeys SET expired = TRUE WHERE owner = $1 RETURNING *;

-- name: ExpireApiKeysByService :many
UPDATE apikeys SET expired = TRUE WHERE service = $1 RETURNING *;

-- name: ExpireExpiredApiKeys :many
UPDATE apikeys SET expired = TRUE WHERE expires_at < $1 RETURNING *;

-- name: DeleteApiKey :exec
DELETE FROM apikeys WHERE api_key = $1;

-- name: DeleteApiKeysByOwner :exec
DELETE FROM apikeys WHERE owner = $1;

-- name: DeleteApiKeysByService :exec
DELETE FROM apikeys WHERE service = $1;

-- name: DeleteExpiredApiKeys :exec
DELETE FROM apikeys WHERE expired = TRUE;
