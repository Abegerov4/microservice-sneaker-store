CREATE TABLE IF NOT EXISTS orders (
    id               TEXT PRIMARY KEY,
    user_id          TEXT NOT NULL,
    total_amount     NUMERIC(10,2) NOT NULL DEFAULT 0,
    status           TEXT NOT NULL DEFAULT 'pending',
    shipping_address TEXT NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS order_items (
    id           TEXT PRIMARY KEY,
    order_id     TEXT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id   TEXT NOT NULL,
    product_name TEXT NOT NULL DEFAULT '',
    quantity     INT NOT NULL DEFAULT 1,
    price        NUMERIC(10,2) NOT NULL DEFAULT 0,
    size         TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
