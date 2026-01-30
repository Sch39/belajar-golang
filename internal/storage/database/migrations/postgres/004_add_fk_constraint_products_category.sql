-- 004_add_fk_constraint_products_category.sql
ALTER TABLE products
ADD CONSTRAINT fk_products_category
FOREIGN KEY (category_id)
REFERENCES categories(id)
ON DELETE SET NULL;