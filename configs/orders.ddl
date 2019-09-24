CREATE SCHEMA IF NOT EXISTS orders AUTHORIZATION lapkin;

CREATE TABLE orders.order_status_code (
	id INT SERIAL PRIMARY KEY,
	status_desc TEXT NOT NULL
);

CREATE TABLE orders.shipping_method (
	id INT SERIAL PRIMARY KEY,
	status_desc TEXT NOT NULL
);

CREATE TABLE orders.payment_method (
	id INT SERIAL PRIMARY KEY,
	method_description TEXT NOT NULL
);

CREATE TABLE orders.customer_payment_method (
	id INT SERIAL PRIMARY KEY,
	customer_id INT REFERENCES customers.customer(id),
	payment_code_id INT REFERENCES orders.payment_method(id)
);

CREATE TABLE orders.order (
	id INT SERIAL PRIMARY KEY,
	order_status_code INT REFERENCES orders.order_status_code(id),
	shipping_status_code INT REFERENCES orders.shipping_status_code(id),
	customer_id INT REFERENCES customers.customer(id),
	order_total INT NOT NULL,
	created_at timestamptz DEFAULT current_timestamp,
	updated_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE orders.order_detail (
	id INT SERIAL PRIMARY KEY,
	order_id INT REFERENCES orders.order(id),
	product_id INT REFERENCES products.product_variant(id),
	quantity INT NOT NULL,
	product_price TEXT NOT NULL,
	created_at timestamptz DEFAULT current_timestamp,
	updated_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE orders.invoice_status_code (
	id INT SERIAL PRIMARY KEY,
	status_desc TEXT NOT NULL
);

CREATE TABLE orders.invoice (
	id INT SERIAL PRIMARY KEY,
	order_id INT REFERENCES orders.order(id),
	status_code INT REFERENCES orders.invoice_status_code(id),
	date timestamptz DEFAULT current_timestamp,
	details TEXT
);