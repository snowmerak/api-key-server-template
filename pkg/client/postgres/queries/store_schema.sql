CREATE TABLE IF NOT EXISTS ApiKeys (
    api_key VARCHAR(512) PRIMARY KEY,
    owner VARCHAR(128) NOT NULL,
    service VARCHAR(128) NOT NULL,
    permissions TEXT NOT NULL,
    payload BYTEA NOT NULL,
    expired BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_api_keys_owner ON ApiKeys (owner);
CREATE INDEX IF NOT EXISTS idx_api_keys_service ON ApiKeys (service);
CREATE INDEX IF NOT EXISTS idx_api_keys_expired ON ApiKeys (expired);
