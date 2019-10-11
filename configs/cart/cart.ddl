CREATE SCHEMA IF NOT EXISTS cart AUTHORIZATION lapkin;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE cart.cart (
	id SERIAL PRIMARY KEY,
	session uuid DEFAULT NULL ,
	customer_id INT REFERENCES customers.customer(id) DEFAULT NULL,
	variant_id INT REFERENCES products.product_variant(id) NOT NULL,
	quantity INT NOT NULL DEFAULT 0
)

uuid_generate_v4 ()