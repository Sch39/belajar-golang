CREATE TABLE IF NOT EXISTS transaction_details (
    id UUID PRIMARY KEY,
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE RESTRICT,
    product_id UUID NOT NULL REFERENCES products(id),

    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price BIGINT NOT NULL CHECK (price > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMPTZ,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CHECK (
        (is_active = true AND deleted_at IS NULL) OR
        (is_active = false AND deleted_at IS NOT NULL)
    )
);

CREATE INDEX idx_transaction_details_product_created
ON transaction_details (product_id, created_at)
WHERE is_active = true;