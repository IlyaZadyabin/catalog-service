CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(256) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add category_id to products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id);

-- Insert initial categories
INSERT INTO categories (code, name) VALUES
('clothing', 'Clothing'),
('shoes', 'Shoes'),
('accessories', 'Accessories');

-- Update products with their categories
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'clothing') 
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'shoes') 
WHERE code IN ('PROD002', 'PROD006');

UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'accessories') 
WHERE code IN ('PROD003', 'PROD005', 'PROD008');
