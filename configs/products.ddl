CREATE SCHEMA IF NOT EXISTS products AUTHORIZATION lapkin;

-- TRIGGERS ********************************************

CREATE OR REPLACE FUNCTION products.update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

CREATE TRIGGER update_product_modtime BEFORE UPDATE ON products.product FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

-- PRODUCT ATTRIBUTE ********************************************

CREATE TABLE products.product_attribute (
	id SERIAL PRIMARY KEY,
	name TEXT,
	display TEXT
);

CREATE TABLE products.product_class (
	id SERIAL PRIMARY KEY,
	name TEXT,
	has_variant BOOLEAN,
	is_shipping_required BOOLEAN
);

CREATE TABLE products.attribute_choice_value (
	id SERIAL PRIMARY KEY,
	display TEXT,
	attribute_id INT REFERENCES products.product_attribute(id)
);

-- CATEGORY ********************************************

CREATE TABLE products.category (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT,
	hidden BOOLEAN,
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
	name TEXT,
	description TEXT,
	price INT,
	available_on BOOLEAN,
	updated_at timestamptz DEFAULT current_timestamp,
	product_class_id INT REFERENCES product_class(id)
);

CREATE TRIGGER update_product_modtime BEFORE UPDATE ON products.product FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE products.product_variant ( 
	id SERIAL PRIMARY KEY,
	sku INT,
	name TEXT,
	price_override INT,
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
	name TEXT
);

CREATE TABLE products.stock (
	id SERIAL PRIMARY KEY,
	qty INT,
	cost_price INT,
	variant_id INT REFERENCES products.product_variant(id),
	quantity_allocated INT,
	location_id INT REFERENCES products.stock_location(id)
);

-- SIZES ********************************************

CREATE TABLE products."size" (
	id serial PRIMARY KEY,
	x TEXT,
	y TEXT
);

	id serial PRIMARY KEY,
	product_class_id INT REFERENCES products.product_class(id),
	product_size_id INT REFERENCES products."size"(id)
);

CREATE TABLE products.product_class_variant_size (
	id serial PRIMARY KEY,
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