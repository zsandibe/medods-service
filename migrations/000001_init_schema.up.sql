CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS sessions (
    id UUID DEFAULT uuid_generate_v4() NOT NULL,
    guid UUID NOT NULL,
    refresh_token BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP
);


CREATE INDEX idx_session_id ON sessions(id);