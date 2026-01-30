-- 003_add_category_id_to_products.sql
ALTER TABLE products 
ADD COLUMN category_id UUID;
