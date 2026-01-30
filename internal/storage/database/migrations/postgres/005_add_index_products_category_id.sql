-- 005_add_index_products_category_id.sql

CREATE INDEX IF NOT EXISTS idx_products_category_id
ON products(category_id);
