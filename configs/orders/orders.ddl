DROP TABLE IF EXISTS "customer_payment_method" CASCADE;
DROP TABLE IF EXISTS "payment_method" CASCADE;
DROP TABLE IF EXISTS "order_status" CASCADE;
DROP TABLE IF EXISTS "order" CASCADE;
DROP TABLE IF EXISTS "shipping_status" CASCADE;
DROP TABLE IF EXISTS "order_detail" CASCADE;
DROP TABLE IF EXISTS "invoice" CASCADE;
DROP TABLE IF EXISTS "invoice_status" CASCADE;

DROP SEQUENCE IF EXISTS "order_status_id_seq" CASCADE;
DROP SEQUENCE IF EXISTS "shipping_status_id_seq" CASCADE;
DROP SEQUENCE IF EXISTS "payment_method_id_seq" CASCADE;
DROP SEQUENCE IF EXISTS "customer_payment_method_id_seq" CASCADE;
DROP SEQUENCE IF EXISTS "order_id_seq" CASCADE;
DROP SEQUENCE IF EXISTS "order_detail_id_seq" CASCADE;
DROP SEQUENCE IF EXISTS "invoice_status_id_seq" CASCADE;
DROP SEQUENCE IF EXISTS "invoice_id_seq" CASCADE;

CREATE SCHEMA IF NOT EXISTS orders AUTHORIZATION lapkin;

CREATE TABLE orders.order_status (
	id SERIAL PRIMARY KEY,
	status TEXT NOT NULL
);

CREATE TABLE orders.shipping_status (
	id SERIAL PRIMARY KEY,
	status TEXT NOT NULL
);

CREATE TABLE orders.payment_method (
	id SERIAL PRIMARY KEY,
	status TEXT NOT NULL
);

CREATE TABLE orders.customer_payment_method (
	id SERIAL PRIMARY KEY,
	customer_id INT REFERENCES customers.customer(id),
	payment_code_id INT REFERENCES orders.payment_method(id)
);

CREATE TABLE orders.order (
	id SERIAL PRIMARY KEY,
	order_status_id INT REFERENCES orders.order_status(id),
	shipping_status_id INT REFERENCES orders.shipping_status(id),
	customer_id INT REFERENCES customers.customer(id),
	order_total INT NOT NULL,
	created_at timestamptz DEFAULT current_timestamp,
	updated_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE orders.order_detail (
	id SERIAL PRIMARY KEY,
	order_id INT REFERENCES orders.order(id),
	variant_id INT REFERENCES products.product_variant(id),
	quantity INT NOT NULL,
	product_price TEXT NOT NULL,
	created_at timestamptz DEFAULT current_timestamp,
	updated_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE orders.invoice_status (
	id SERIAL PRIMARY KEY,
	status TEXT NOT NULL
);

CREATE TABLE orders.invoice (
	id SERIAL PRIMARY KEY,
	order_id INT REFERENCES orders.order(id),
	invoice_status_id INT REFERENCES orders.invoice_status(id),
	date timestamptz DEFAULT current_timestamp,
	details TEXT
);
