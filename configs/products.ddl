DROP TABLE if EXISTS "product_attribute" cascade;
DROP TABLE if EXISTS "attribute_choice_value" cascade;
DROP TABLE if EXISTS "category" cascade;
DROP TABLE if EXISTS "product_class" cascade;
DROP TABLE if EXISTS "product_class_product_attribute" cascade;
DROP TABLE if EXISTS "product_class_variant_attribute" cascade;
DROP TABLE if EXISTS "product" cascade;
DROP TABLE if EXISTS "product_variant" cascade;
DROP TABLE if EXISTS "product_categories" cascade;
DROP TABLE if EXISTS "product_images" cascade;
DROP TABLE if EXISTS "variant_images" cascade;
DROP TABLE if EXISTS "stock" cascade;
DROP TABLE if EXISTS "stock_location" cascade;
DROP TABLE if EXISTS "product_class_product_size" cascade;
DROP TABLE if EXISTS "size" cascade;
DROP TABLE if EXISTS "product_variant_to_size" cascade;

DROP SEQUENCE if EXISTS "product_attribute_id_seq" cascade;
DROP SEQUENCE if EXISTS "product_class_id_seq" cascade;
DROP SEQUENCE if EXISTS "attribute_choice_value_id_seq" cascade;
DROP SEQUENCE if EXISTS "category_id_seq" cascade;
DROP SEQUENCE if EXISTS "product_class_product_attribute_id_seq" cascade;
DROP SEQUENCE if EXISTS "product_class_variant_attribute_id_seq" cascade;
DROP SEQUENCE if EXISTS "product_id_seq" cascade;
DROP SEQUENCE if EXISTS "product_variant_id_seq" cascade;
DROP SEQUENCE if EXISTS "product_categories_id_seq" cascade;
DROP SEQUENCE if EXISTS "product_images_id_seq" cascade;
DROP SEQUENCE if EXISTS "variant_image_id_seq" cascade;
DROP SEQUENCE if EXISTS "stock_location_id_seq" cascade;
DROP SEQUENCE if EXISTS "stock_id_seq" cascade;
DROP SEQUENCE if EXISTS "size_id_seq" cascade;
DROP SEQUENCE if EXISTS "product_class_product_size_id_seq" cascade;
DROP SEQUENCE if EXISTS "product_variant_to_size_id_seq" cascade;

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

-- Functions ********************************************

-- get_products
DROP FUNCTION products.get_products(INT);
CREATE OR REPLACE FUNCTION products.get_products(INT)
RETURNS TABLE  (
	id INT,
	name TEXT,
	description TEXT,
	price INT,
	size TEXT
)
AS $$
	BEGIN
	 	RETURN QUERY 
 		SELECT 
 			p.id, 
 			p."name", 
 			p.description, 
 			p.price,
 			ss.x || 'x' || ss.y
		FROM 
			products.product p
		INNER JOIN 
			products.product_categories pc ON pc.product_id = p.id AND pc.category_id = $1
		INNER JOIN
			products.product_variant pv ON pv.product_id = p.id AND pv."attributes" @> '{"parent": true}'
		INNER JOIN
			products.product_variant_to_size pvts ON pvts.variant_id = pv.id
		INNER JOIN
			products."size" ss ON ss.id = pvts.product_size_id
		GROUP BY p.id, ss.x, ss.y;
	END;
$$ LANGUAGE plpgsql;


SELECT * FROM products.get_products(8);
	
DROP FUNCTION products.get_variant(INT, TEXT);
CREATE OR REPLACE FUNCTION products.get_variant(p_id INT, p_size TEXT)
RETURNS TABLE (
	variant_id INT,
	product_id INT,
	name TEXT,
	description TEXT,
	price_override INT,
	ATTRIBUTES jsonb,
	sizes TEXT[],
	"size" TEXT,
	images TEXT[]
)
AS $$
	BEGIN
	 	RETURN QUERY
	 	SELECT 
			pv.id AS variant_id,
			p.id AS product_id,
			pv."name",
			p.description,
			pv.price_override, 
			pv."attributes",
			array_agg(DISTINCT s.x || 'x' || s.y) AS sizes,
			ss.x || 'x' || ss.y,
			array_agg(DISTINCT p."name" || '/' || ss.x || 'x' || ss.y || '/' || vim.name) AS images
		FROM
			products.product_variant pv
		INNER JOIN
			products.product p ON pv.product_id = p.id
		INNER JOIN
			products.product_class pc ON pc.id = p.product_class_id
		INNER JOIN
			products.product_class_product_size pcps ON pcps.product_class_id = pc.id
		INNER JOIN
			products.product_variant_to_size pvts ON pvts.variant_id = pv.id
		INNER JOIN
			products."size" s ON s.id = pcps.product_size_id
		INNER JOIN
			products."size" ss ON ss.id = pvts.product_size_id
		INNER JOIN 
			products.variant_images vim ON vim.variant_id = pv.id
		WHERE 
			p.id = p_id
		AND 
			(length(p_size) = 0 OR ss.x || 'x' || ss.y = p_size)
		GROUP BY
			p.id, pv.id, ss.x, ss.y
		ORDER BY p.id;
	END;
$$ LANGUAGE plpgsql;

SELECT * FROM products.get_variant(3, '');

DROP FUNCTION IF EXISTS products.get_categories(INT);
CREATE OR REPLACE FUNCTION products.get_categories(category_parent_id int)
RETURNS TABLE (
	categories json
)
AS $$
	BEGIN
	 	RETURN QUERY 
		SELECT array_to_json(array_agg(a.category::json))
		FROM
			(
				SELECT
					json_build_object('category_name', t.display, 'sub_categories', json_agg(t.categories)) AS category
				FROM (
					SELECT
						pcl.display AS display,
						json_build_object('display', c.display, 'url', c."name") AS categories
					FROM
						products.category c
					INNER JOIN
						products.product_categories pc ON pc.category_id = c.id
					INNER JOIN
						products.product p ON p.id = pc.product_id
					INNER JOIN
						products.product_class pcl ON pcl.id = p.product_class_id
					WHERE
						c.parent_id = category_parent_id
					AND
						c.hidden = FALSE
					GROUP BY 
						pcl.display, c.display, c."name"
				) t
				GROUP BY
					t.display
			) a;
	END;
$$ LANGUAGE plpgsql;

SELECT * FROM products.get_categories(1);

-- TRIGGERS ********************************************

DROP TRIGGER IF EXISTS "update_modtime" ON products.product;
DROP TRIGGER IF EXISTS "update_modtime" ON products.product_variant;
DROP TRIGGER IF EXISTS "add_sku_tr" ON products.product_variant;

CREATE OR REPLACE FUNCTION products.update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

CREATE TRIGGER update_modtime BEFORE UPDATE ON products.product FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER update_modtime BEFORE UPDATE ON products.product_variant FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

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

SELECT array_agg(a.*)
FROM
	(
		SELECT
			json_build_object('category', t.display, 'categories', json_agg(t.categories)) AS result
		FROM (
			SELECT
				c.display AS display,
				json_build_object('display', c.display, 'url', c."name") AS categories
			FROM
				products.category c
			INNER JOIN
				products.product_categories pc ON pc.category_id = c.id
			INNER JOIN
				products.product p ON p.id = pc.product_id
			INNER JOIN
				products.product_class pcl ON pcl.id = p.product_class_id
			WHERE
				c.parent_id = 1
			GROUP BY 
				pcl.display, c.display, c."name"
		) t
		GROUP BY
			t.display
	) a
	

