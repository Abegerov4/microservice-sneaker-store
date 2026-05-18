SELECT 'CREATE DATABASE sneaker_products' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'sneaker_products')\gexec
SELECT 'CREATE DATABASE sneaker_orders' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'sneaker_orders')\gexec
SELECT 'CREATE DATABASE sneaker_users' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'sneaker_users')\gexec
SELECT 'CREATE DATABASE sneaker_ai' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'sneaker_ai')\gexec
