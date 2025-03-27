-- +goose Up
CREATE TABLE IF NOT EXISTS metrics (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL CHECK (type IN ('gauge', 'counter')),
    value DOUBLE PRECISION,
    delta BIGINT
);

-- +goose Down
DROP TABLE IF EXISTS metrics;
