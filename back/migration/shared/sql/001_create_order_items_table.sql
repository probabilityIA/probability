-- Create order_items table (PostgreSQL)
CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    order_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(64),

    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(12,2) NOT NULL,
    total_price DECIMAL(12,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',

    discount DECIMAL(12,2) DEFAULT 0,
    discount_percent DECIMAL(5,2) DEFAULT 0,
    tax DECIMAL(12,2) DEFAULT 0,
    tax_rate DECIMAL(5,4),

    unit_price_base DECIMAL(12,2) NOT NULL DEFAULT 0,
    unit_price_base_presentment DECIMAL(12,2) NOT NULL DEFAULT 0,
    unit_price_presentment DECIMAL(12,2) NOT NULL DEFAULT 0,
    total_price_presentment DECIMAL(12,2) NOT NULL DEFAULT 0,
    discount_presentment DECIMAL(12,2) DEFAULT 0,
    tax_presentment DECIMAL(12,2) DEFAULT 0,

    product_sku VARCHAR(255),
    product_name VARCHAR(255),
    variant_id VARCHAR(255),
    variant_label VARCHAR(255),
    fulfillment_status VARCHAR(64),
    metadata JSONB,

    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);
CREATE INDEX IF NOT EXISTS idx_order_items_deleted_at ON order_items(deleted_at);
