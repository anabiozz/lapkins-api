CREATE SCHEMA IF NOT EXISTS orders AUTHORIZATION lapkin;

CREATE TABLE orders.order_status_code (
	id INT SERIAL PRIMARY KEY,
	status_desc TEXT
)

CREATE TABLE orders.order (
	id INT SERIAL PRIMARY KEY,
	status_code INT REFERENCES orders.order_status_code(id),
	customer_id INT REFERENCES customers.customer(id),
	order_total TEXT,
	created_at timestamptz DEFAULT current_timestamp,
	updated_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE orders.order_detail (
	id INT SERIAL PRIMARY KEY,
	order_id INT REFERENCES orders.order(id),
	product_id INT REFERENCES products.product_variant(id),
	quantity INT,
	product_price TEXT,
	created_at timestamptz DEFAULT current_timestamp,
	updated_at timestamptz DEFAULT current_timestamp
);