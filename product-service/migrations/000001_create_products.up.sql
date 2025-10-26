CREATE TABLE IF NOT EXISTS products (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    brand       TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    price       NUMERIC(10,2) NOT NULL DEFAULT 0,
    sizes       TEXT[] NOT NULL DEFAULT '{}',
    stock       INT NOT NULL DEFAULT 0,
    image_url   TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand);
CREATE INDEX IF NOT EXISTS idx_products_price ON products(price);
