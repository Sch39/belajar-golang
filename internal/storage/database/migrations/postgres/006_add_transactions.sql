CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    total_price BIGINT NOT NULL CHECK (total_price > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMPTZ,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CHECK (
        (is_active = true AND deleted_at IS NULL) OR
        (is_active = false AND deleted_at IS NOT NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_transactions_id ON transactions(id)
WHERE is_active = true;