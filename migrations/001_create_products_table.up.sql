CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    store_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    amount INTEGER NOT NULL DEFAULT 0,
    price NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_store_id ON products(store_id);
CREATE INDEX idx_products_created_at ON products(created_at);