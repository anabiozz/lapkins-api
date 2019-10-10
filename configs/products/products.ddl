DROP TABLE IF EXISTS "product_attribute" cascade;
DROP TABLE IF EXISTS "attribute_choice_value" cascade;
DROP TABLE IF EXISTS "category" cascade;
DROP TABLE IF EXISTS "product_class" cascade;
DROP TABLE IF EXISTS "product_class_product_attribute" cascade;
DROP TABLE IF EXISTS "product_class_variant_attribute" cascade;
DROP TABLE IF EXISTS "product" cascade;
DROP TABLE IF EXISTS "product_variant" cascade;
DROP TABLE IF EXISTS "product_categories" cascade;
DROP TABLE IF EXISTS "product_images" cascade;
DROP TABLE IF EXISTS "variant_images" cascade;
DROP TABLE IF EXISTS "stock" cascade;
DROP TABLE IF EXISTS "stock_location" cascade;
DROP TABLE IF EXISTS "product_class_product_size" cascade;
DROP TABLE IF EXISTS "size" cascade;
DROP TABLE IF EXISTS "product_variant_to_size" cascade;

DROP SEQUENCE IF EXISTS "product_attribute_id_seq" cascade;
DROP SEQUENCE IF EXISTS "product_class_id_seq" cascade;
DROP SEQUENCE IF EXISTS "attribute_choice_value_id_seq" cascade;
DROP SEQUENCE IF EXISTS "category_id_seq" cascade;
DROP SEQUENCE IF EXISTS "product_class_product_attribute_id_seq" cascade;
DROP SEQUENCE IF EXISTS "product_class_variant_attribute_id_seq" cascade;
DROP SEQUENCE IF EXISTS "product_id_seq" cascade;
DROP SEQUENCE IF EXISTS "product_variant_id_seq" cascade;
DROP SEQUENCE IF EXISTS "product_categories_id_seq" cascade;
DROP SEQUENCE IF EXISTS "product_images_id_seq" cascade;
DROP SEQUENCE IF EXISTS "variant_image_id_seq" cascade;
DROP SEQUENCE IF EXISTS "stock_location_id_seq" cascade;
DROP SEQUENCE IF EXISTS "stock_id_seq" cascade;
DROP SEQUENCE IF EXISTS "size_id_seq" cascade;
DROP SEQUENCE IF EXISTS "product_class_product_size_id_seq" cascade;
DROP SEQUENCE IF EXISTS "product_variant_to_size_id_seq" cascade;

CREATE SCHEMA IF NOT EXISTS products AUTHORIZATION lapkin;

-- PRODUCT ATTRIBUTE ********************************************

CREATE TABLE products.product_attribute (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	display TEXT NOT NULL,
	created_at timestamp with time zone DEFAULT current_timestamp
);

CREATE TABLE products.product_class (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	display TEXT NOT NULL,
	sku_part TEXT UNIQUE,
	created_at timestamp with time zone DEFAULT current_timestamp
);

CREATE TABLE products.attribute_choice_value (
	id SERIAL PRIMARY KEY,
	display TEXT NOT NULL,
	attribute_id INT REFERENCES products.product_attribute(id),
	created_at timestamp with time zone DEFAULT current_timestamp
);

-- CATEGORY ********************************************

CREATE TABLE products.category (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	display TEXT NOT NULL,
	description TEXT,
	hidden BOOLEAN NOT NULL DEFAULT false,
	tree_id INT,
	parent_id INT REFERENCES products.category(id),
	created_at timestamp with time zone DEFAULT current_timestamp
);

CREATE TABLE products.product_class_product_attribute (
	id SERIAL PRIMARY KEY,
	product_class_id INT REFERENCES products.product_class(id),
	product_attribute_id INT REFERENCES products.product_attribute(id),
	created_at timestamp with time zone DEFAULT current_timestamp
);

CREATE TABLE products.product_class_variant_attribute (
	id SERIAL PRIMARY KEY,
	product_class_id INT REFERENCES products.product_class(id),
	product_attribute_id INT REFERENCES products.product_attribute(id),
	created_at timestamp with time zone DEFAULT current_timestamp
);

-- PRODUCT ********************************************

CREATE TABLE products.product (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	price INT NOT NULL,
	product_class_id INT REFERENCES product_class(id),
	created_at timestamptz DEFAULT current_timestamp,
	updated_at timestamptz DEFAULT current_timestamp
);


CREATE TABLE products.product_variant ( 
	id SERIAL PRIMARY KEY,
	sku TEXT,
	name TEXT NOT NULL,
	price_override INT NOT NULL,
	product_id INT REFERENCES products.product(id),
--	parent_id INT REFERENCES products.product_variant(id),
	attributes JSONB,
	created_at timestamptz DEFAULT current_timestamp,
	updated_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE products.product_categories (
	id SERIAL PRIMARY KEY,
	product_id INT REFERENCES products.product(id),
	category_id INT REFERENCES products.category(id),
	created_at timestamptz DEFAULT current_timestamp
);

-- IMAGES ********************************************

--CREATE TABLE products.product_images (
--	id SERIAL PRIMARY KEY,
--	image TEXT,
--	alt TEXT,
--	"order" INT,
--	product_id INT REFERENCES products.product(id),
--	created_at timestamptz DEFAULT current_timestamp
--);

CREATE TABLE products.variant_images (
	id SERIAL PRIMARY KEY,
	name TEXT,
	variant_id INT REFERENCES products.product_variant(id),
	created_at timestamptz DEFAULT current_timestamp
);

-- STOCK ********************************************

CREATE TABLE products.stock_location (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	created_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE products.stock (
	id SERIAL PRIMARY KEY,
	qty INT NOT NULL,
	cost_price INT NOT NULL,
	variant_id INT REFERENCES products.product_variant(id),
	location_id INT REFERENCES products.stock_location(id),
	created_at timestamptz DEFAULT current_timestamp
);

-- SIZES ********************************************

CREATE TABLE products."size" (
	id SERIAL PRIMARY KEY,
	x TEXT NOT NULL,
	y TEXT NOT NULL,
	created_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE products.product_class_product_size(
	id SERIAL PRIMARY KEY,
	product_class_id INT REFERENCES products.product_class(id),
	product_size_id INT REFERENCES products."size"(id),
	created_at timestamptz DEFAULT current_timestamp
);

CREATE TABLE products.product_variant_to_size (
	id SERIAL PRIMARY KEY,
	variant_id INT REFERENCES products.product_variant(id),
	product_size_id INT REFERENCES products."size"(id),
	created_at timestamptz DEFAULT current_timestamp
);