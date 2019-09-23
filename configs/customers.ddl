CREATE SCHEMA IF NOT EXISTS customers AUTHORIZATION lapkin;

CREATE TABLE customers.customer (
	id SERIAL PRIMARY KEY,
	email TEXT,
	name TEXT,
	gender TEXT,
	password TEXT,
	first_name TEXT,
	surname TEXT,
	date_ofbirth DATE,
	phone TEXT,
	created_at timestamptz DEFAULT current_timestamp,
	updated_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE customers.payment_method (
	id INT SERIAL PRIMARY KEY,
	method_description TEXT
)

CREATE TABLE customers.customer_payment_method (
	id INT SERIAL PRIMARY KEY,
	customer_id INT REFERENCES customers.customer(id),
	payment_code_id INT REFERENCES customers.payment_method(id)
)