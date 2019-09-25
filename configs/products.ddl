drop table if exists "product_attribute" cascade;
drop table if exists "attribute_choice_value" cascade;
drop table if exists "category" cascade;
drop table if exists "product_class" cascade;
drop table if exists "product_class_product_attribute" cascade;
drop table if exists "product_class_variant_attribute" cascade;
drop table if exists "product" cascade;
drop table if exists "product_variant" cascade;
drop table if exists "product_categories" cascade;
drop table if exists "product_images" cascade;
drop table if exists "variant_image" cascade;
drop table if exists "stock" cascade;
drop table if exists "stock_location" cascade;
drop table if exists "product_class_product_size" cascade;
drop table if exists "size" cascade;
drop table if exists "product_variant_to_size" cascade;

drop sequence if exists "product_attribute_id_seq" cascade;
drop sequence if exists "product_class_id_seq" cascade;
drop sequence if exists "attribute_choice_value_id_seq" cascade;
drop sequence if exists "category_id_seq" cascade;
drop sequence if exists "product_class_product_attribute_id_seq" cascade;
drop sequence if exists "product_class_variant_attribute_id_seq" cascade;
drop sequence if exists "product_id_seq" cascade;
drop sequence if exists "product_variant_id_seq" cascade;
drop sequence if exists "product_categories_id_seq" cascade;
drop sequence if exists "product_images_id_seq" cascade;
drop sequence if exists "variant_image_id_seq" cascade;
drop sequence if exists "stock_location_id_seq" cascade;
drop sequence if exists "stock_id_seq" cascade;
drop sequence if exists "size_id_seq" cascade;
drop sequence if exists "product_class_product_size_id_seq" cascade;
drop sequence if exists "product_variant_to_size_id_seq" cascade;

DROP TRIGGER IF EXISTS "update_modified_column" ON products.product;
DROP TRIGGER IF EXISTS "add_sku_tr" ON products.product_variant;

CREATE SCHEMA IF NOT EXISTS products AUTHORIZATION lapkin;

-- PRODUCT ATTRIBUTE ********************************************

CREATE TABLE products.product_attribute (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	display TEXT NOT NULL
);

CREATE TABLE products.product_class (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	sku_part TEXT UNIQUE
);

CREATE TABLE products.attribute_choice_value (
	id SERIAL PRIMARY KEY,
	display TEXT NOT NULL,
	attribute_id INT REFERENCES products.product_attribute(id)
);

-- CATEGORY ********************************************

CREATE TABLE products.category (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT,
	hidden BOOLEAN NOT NULL DEFAULT false,
	tree_id INT,
	parent_id INT REFERENCES products.category(id)
);

CREATE TABLE products.product_class_product_attribute (
	id SERIAL PRIMARY KEY,
	product_class_id INT REFERENCES products.product_class(id),
	product_attribute_id INT REFERENCES products.product_attribute(id)
);

CREATE TABLE products.product_class_variant_attribute (
	id SERIAL PRIMARY KEY,
	product_class_id INT REFERENCES products.product_class(id),
	product_attribute_id INT REFERENCES products.product_attribute(id)
);

-- PRODUCT ********************************************

CREATE TABLE products.product (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	price INT NOT NULL,
	updated_at timestamptz DEFAULT current_timestamp NOT NULL,
	product_class_id INT REFERENCES product_class(id)
);


CREATE TABLE products.product_variant ( 
	id SERIAL PRIMARY KEY,
	sku TEXT,
	name TEXT NOT NULL,
	price_override INT NOT NULL,
	product_id INT REFERENCES products.product(id),
	attributes JSONB
);

CREATE TABLE products.product_categories (
	id SERIAL PRIMARY KEY,
	product_id INT REFERENCES products.product(id),
	category_id INT REFERENCES products.category(id)
);

-- IMAGES ********************************************

CREATE TABLE products.product_images (
	id SERIAL PRIMARY KEY,
	image TEXT,
	ppoi TEXT,
	alt TEXT,
	"order" INT,
	product_id INT REFERENCES products.product(id)
);

CREATE TABLE products.variant_image (
	id SERIAL PRIMARY KEY,
	image_id INT REFERENCES products.product_images(id),
	variant_id INT REFERENCES products.product_variant(id)
);

-- STOCK ********************************************

CREATE TABLE products.stock_location (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE products.stock (
	id SERIAL PRIMARY KEY,
	qty INT NOT NULL,
	cost_price INT NOT NULL,
	variant_id INT REFERENCES products.product_variant(id),
	location_id INT REFERENCES products.stock_location(id)
);

-- SIZES ********************************************

CREATE TABLE products."size" (
	id SERIAL PRIMARY KEY,
	x TEXT NOT NULL,
	y TEXT NOT NULL
);

CREATE TABLE products.product_class_product_size(
	id SERIAL PRIMARY KEY,
	product_class_id INT REFERENCES products.product_class(id),
	product_size_id INT REFERENCES products."size"(id)
);

CREATE TABLE products.product_variant_to_size (
	id SERIAL PRIMARY KEY,
	variant_id INT REFERENCES products.product_variant(id),
	product_size_id INT REFERENCES products."size"(id)
);

-- Functions ********************************************

-- get_products
CREATE OR REPLACE FUNCTION products.get_products(INT)
RETURNS TABLE  (
	id INT,
	name TEXT,
	description TEXT,
	price INT
)
AS $$
	BEGIN
	 	RETURN QUERY SELECT product.id, product."name", product.description, product.price
		FROM products.product
		JOIN products.product_categories AS pc ON pc.product_id = product.id AND pc.category_id = $1;
	END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION products.get_product_variant_by_id(p_id INT, p_size TEXT)
RETURNS TABLE (
	id INT,
	product_id INT,
	name TEXT,
	description TEXT,
	price_override INT,
	ATTRIBUTES jsonb,
	sizes TEXT[]
)
AS $$
	BEGIN
	 	RETURN QUERY 
 		SELECT 
			pv.id,
			p.id,
			pv."name",
			p.description,
			pv.price_override, 
			pv."attributes",
			ARRAY_AGG(s.x || 'x' || s.y) AS size
		FROM products.product_variant AS pv
		JOIN products.product AS p ON pv.product_id = p.id
		JOIN products.product_class_product_size AS pcps ON pcps.product_class_id = p.product_class_id
		JOIN products."size" AS s ON s.id = pcps.product_size_id
		WHERE 
			string_to_array(COALESCE(NULLIF(p_size, ''), (
												SELECT s.x || 'x' || s.y AS size 
												FROM products.product_class_variant_size AS pcvs
												JOIN products."size" AS s ON s.id = pcvs.product_size_id
												WHERE pcvs.variant_id = pv.id
											)
									), '|') IN 
				(
					SELECT ARRAY_AGG(s.x || 'x' || s.y) AS size
					FROM products.product_class_variant_size AS pcvs
					JOIN products."size" AS s ON s.id = pcvs.product_size_id
					WHERE pcvs.variant_id = pv.id
					
				)
		AND pv."attributes"->'frame' IS NULL 
		AND pv.product_id = p_id
		GROUP BY pv.id, p.id;
	END;
$$ LANGUAGE plpgsql;

-- TRIGGERS ********************************************

CREATE OR REPLACE FUNCTION products.update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

CREATE TRIGGER update_product_modtime BEFORE UPDATE ON products.product FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

--currval(pg_get_serial_sequence(TG_TABLE_NAME, 'id'))

CREATE OR REPLACE FUNCTION products.add_sku()
RETURNS TRIGGER AS $$
BEGIN
	NEW.sku = (
    	SELECT pc.sku_part || '-' || currval(pg_get_serial_sequence(TG_TABLE_NAME, 'id'))
	   	FROM products.product_class AS pc
	   	JOIN products.product AS p ON p.product_class_id = pc.id
	   	JOIN (SELECT NEW.*) AS pv ON pv.product_id = p.id AND pv.id = currval(pg_get_serial_sequence(TG_TABLE_NAME, 'id'))
	);
    RETURN NEW;   
END;
$$ language 'plpgsql';

CREATE TRIGGER add_sku_tr BEFORE INSERT ON products.product_variant FOR EACH ROW EXECUTE PROCEDURE products.add_sku();