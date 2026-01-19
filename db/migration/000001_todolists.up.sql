CREATE TABLE IF NOT EXISTS todolists (
    id BIGSERIAL PRIMARY KEY,
    created_by BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- migrate create -ext sql -dir db/migration -seq posts -format
-- migrate -path ./db/migration -database "postgres://ankush:postgres@localhost/blogapi?sslmode=disable" up 1