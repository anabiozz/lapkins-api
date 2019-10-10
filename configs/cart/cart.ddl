CREATE SCHEMA IF NOT EXISTS cart AUTHORIZATION lapkin;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE cart.cart (
	session uuid DEFAULT uuid_generate_v4 (),
	customer_id INT REFERENCES customers.customer(id) DEFAULT NULL,
	variant_id INT REFERENCES products.product_variant(id) NOT NULL,
	quantity INT NOT NULL DEFAULT 0,
	PRIMARY KEY(session)
)

CREATE TABLE cart.cart (
	session uuid DEFAULT uuid_generate_v4 (),
	data JSONB NOT NULL,
	PRIMARY KEY(session)
)