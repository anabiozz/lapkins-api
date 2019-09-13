CREATE SCHEMA IF NOT EXISTS products AUTHORIZATION lapkin;

-- FUNCTIONS ********************************************

CREATE OR REPLACE FUNCTION products.update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';


-- PRODUCT ATTRIBUTE ********************************************

CREATE TABLE products.product_attribute (
	id SERIAL PRIMARY KEY,
	name TEXT,
	display TEXT
);

INSERT INTO products.product_attribute (name, display)
VALUES ('authors', 'авторы'), ('materials', 'материалы'), ('finish', 'покрытие'), ('print type', 'тип печати'), ('postcards size', 'размер'), ('posters size', 'размер');


CREATE TABLE products.product_class (
	id SERIAL PRIMARY KEY,
	name TEXT,
	has_variant BOOLEAN,
	is_shipping_required BOOLEAN
);

INSERT INTO products.product_class (name)
VALUES ('postcards'), ('posters'), ('badges'), ('table lamps');


CREATE TABLE products.attribute_choice_value (
	id SERIAL PRIMARY KEY,
	display TEXT,
	attribute_id INT REFERENCES products.product_attribute(id)
);

INSERT INTO products.attribute_choice_value (display, attribute_id)
VALUES ('Анастасия Кондратьева', 1), 
('Lolka Lolkina', 1),
('240 g/m² pure white paper', 2),
('300 g/m² Munken Lynx Rough paper (woodfree)', 2),
('semi-gloss', 3),
('gloss', 3),
('matte', 3),
('digital printing', 4),
('12-colour digital printing', 4),
('105х148', 5),
('148x210', 5),
('130x180', 5),
('300x450', 5),
('200x300', 5),
('400x600', 5),
('600x900', 5),
('1000x1500', 5),
('800x1200', 5),
('A4', 5),
('A6', 5);


-- CATEGORY ********************************************

CREATE TABLE products.category (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT,
	hidden BOOLEAN,
	lft INT,
	rgt INT,
	tree_id INT,
	parent_id INT REFERENCES products.category(id)
);

INSERT INTO products.category (name)
VALUES ('wallart'), ('stationery'), ('gifts'), ('home');



CREATE TABLE products.product_class_product_attribute (
	id SERIAL PRIMARY KEY,
	product_class_id INT REFERENCES products.product_class(id),
	product_attribute_id INT REFERENCES products.product_attribute(id)
);

INSERT INTO products.product_class_product_attribute (product_class_id, product_attribute_id)
VALUES (1, 1), (1, 3), (1, 4), (1, 5), (2, 1), (2, 3), (2, 4), (2, 6);

CREATE TABLE products.product_class_variant_attribute (
	id SERIAL PRIMARY KEY,
	product_class_id INT REFERENCES products.product_class(id),
	product_attribute_id INT REFERENCES products.product_attribute(id)
);

INSERT INTO products.product_class_variant_attribute (product_class_id, product_attribute_id)
VALUES (1, 2), (2, 2);

-- PRODUCT ********************************************

CREATE TABLE products.product (
	id SERIAL PRIMARY KEY,
	name TEXT,
	description TEXT,
	price INT,
	weight INT,
	sizes TEXT[],
	available_on BOOLEAN,
	updated_at timestamptz DEFAULT current_timestamp,
	product_class_id INT REFERENCES product_class(id)
);

INSERT INTO products.product (name, description, price, product_class_id, sizes)
VALUES ('плакат веселье', 'Et deserunt labore excepteur id eiusmod reprehenderit do nostrud cupidatat consectetur laboris culpa.', 300, 2, '{"130x180", "300x450", "200x300", "400x600", "600x900"}'),
('плакат надпись со смыслом', 'Commodo labore est qui laboris irure esse aliquip', 300, 2, '{"130x180", "300x450", "200x300", "400x600", "600x900"}'),
('открытка веселье', 'Anim ex occaecat occaecat non tempor in enim id mollit.', 50, 1, '{"105х148", "148x210"}'),
('открытка надпись со смыслом', 'Nisi eiusmod laborum ullamco mollit elit amet deserunt ex sit nisi consectetur cillum commodo incididunt.', 50, 1, '{"105х148", "148x210"}');

CREATE TRIGGER update_product_modtime BEFORE UPDATE ON products.product FOR EACH ROW EXECUTE PROCEDURE update_modified_column();

CREATE TABLE products.product_variant (
	id SERIAL PRIMARY KEY,
	sku INT,
	name TEXT,
	price_override INT,
	weight_override INT,
	product_id INT REFERENCES products.product(id),
	attributes JSONB
);

INSERT INTO products.product_variant (name, price_override, product_id, attributes)
VALUES ('плакат веселье 300x450', 300, 1, '{"sizes": ["300x450"], "authors": ["Анастасия Кондратьева", "Lolka Lolkina"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper"]}'), 
('плакат веселье с металической рамой 300x450', 600, 1, '{"sizes": ["300x450"], "authors": ["Анастасия Кондратьева", "Lolka Lolkina"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper", "металл"]}'),
('плакат веселье с деревянной рамой 300x450', 400, 1, '{"sizes": ["300x450"], "authors": ["Анастасия Кондратьева", "Lolka Lolkina"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper", "дерево"]}'),
('плакат веселье 600x900', 300, 1, '{"sizes": ["600x900"], "authors": ["Анастасия Кондратьева", "Lolka Lolkina"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper"]}'), 
('плакат веселье с металической рамой 600x900', 600, 1, '{"sizes": ["600x900"], "authors": ["Анастасия Кондратьева", "Lolka Lolkina"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper", "металл"]}'),
('плакат веселье с деревянной рамой 600x900', 400, 1, '{"sizes": ["600x900"], "authors": ["Анастасия Кондратьева", "Lolka Lolkina"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper", "дерево"]}'),
('открытка веселье 105х148', 50, 3, '{"sizes": ["105х148"], "authors": ["Анастасия Кондратьева", "Lolka Lolkina"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["300 g/m² Munken Lynx Rough paper (woodfree)"]}'),
('плакат надпись со смыслом 300x450', 300, 2, '{"sizes": ["300x450"], "authors": ["Анастасия Кондратьева"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper"]}'), 
('плакат надпись со смыслом с металической рамой 300x450', 600, 2, '{"sizes": ["300x450"], "authors": ["Анастасия Кондратьева"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper", "металл"]}'),
('плакат надпись со смыслом с деревянной рамой 300x450', 400, 2, '{"sizes": ["300x450"], "authors": ["Анастасия Кондратьева"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper", "дерево"]}'),
('плакат надпись со смыслом 600x900', 300, 2, '{"sizes": ["600x900"], "authors": ["Анастасия Кондратьева"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper"]}'), 
('плакат надпись со смыслом с металической рамой 600x900', 600, 2, '{"sizes": ["600x900"], "authors": ["Анастасия Кондратьева"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper", "металл"]}'),
('плакат надпись со смыслом с деревянной рамой 600x900', 400, 2, '{"sizes": ["600x900"], "authors": ["Анастасия Кондратьева"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["240 g/m² pure white paper", "дерево"]}'),
('открытка надпись со смыслом 105х148', 50, 4, '{"sizes": ["105х148"], "authors": ["Анастасия Кондратьева"], "finish": ["semi-gloss"], "print type": ["digital printing"], "materials": ["300 g/m² Munken Lynx Rough paper (woodfree)"]}');



CREATE TABLE products.product_categories (
	id SERIAL PRIMARY KEY,
	product_id INT REFERENCES products.product(id),
	category_id INT REFERENCES products.category(id)
);

INSERT INTO products.product_categories (product_id, category_id)
VALUES (1, 1), (2, 1), (3, 2), (4, 2);

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

CREATE OR REPLACE FUNCTION products.get_product_by_id(INT)
RETURNS TABLE  (
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
	 		product_variant.id,
	 		p.id,
	 		product_variant."name",
	 		p.description,
	 		product_variant.price_override, 
	 		product_variant."attributes",
	 		p.sizes::TEXT[]
	 	FROM products.product_variant
	 	JOIN products.product AS p ON products.product_variant.product_id = p.id
	 	WHERE products.product_variant.product_id = $1;
	END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION products.get_product_variant_by_id(INT)
RETURNS TABLE  (
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
	 		product_variant.id,
	 		p.id,
	 		product_variant."name",
	 		p.description,
	 		product_variant.price_override, 
	 		product_variant."attributes",
	 		p.sizes::TEXT[]
	 	FROM products.product_variant
	 	JOIN products.product AS p ON products.product_variant.product_id = p.id
	 	WHERE products.product_variant.id = $1;
	END;
$$ LANGUAGE plpgsql;

-- **************************************************************************

SELECT *
FROM product_variant  
WHERE product_variant."attributes" ->> '600x900' AND product_id = 2


SELECT * 
FROM product_variant 
WHERE "attributes" @> '{"sizes": ["600x900"]}' AND product_id = 2

UPDATE products.product_variant
SET ATTRIBUTES = (
WITH attr AS (
	SELECT pcva.product_attribute_id AS id
	FROM product_class_variant_attribute AS pcva
	WHERE pcva.product_class_id = 1
	UNION
	SELECT pcpa.product_attribute_id AS id
	FROM product_class_product_attribute AS pcpa
	WHERE pcpa.product_class_id = 1
)
SELECT json_object_agg(ch.name, ch.display)
FROM (
	SELECT p_attr.name, json_agg(acv.display) AS display
	FROM attr
	JOIN product_attribute AS p_attr ON p_attr.id = attr.id
	JOIN attribute_choice_value AS acv ON acv.attribute_id = attr.id
	GROUP BY p_attr.name
) ch)
WHERE product_class_id = 1;


WITH attr AS (
	SELECT * FROM product_attribute
)
SELECT json_object_agg(ch.name, ch.display)
FROM (
	SELECT attr.name, json_agg(acv.display) AS display
	FROM attr
	JOIN attribute_choice_value AS acv ON attr.id = acv.attribute_id
	LEFT OUTER JOIN product_class_variant_attribute AS pcva ON attr.id = pcva.product_attribute_id AND pcva.product_class_id = 2
	LEFT OUTER JOIN product_class_product_attribute AS pcpa ON attr.id = pcpa.product_attribute_id AND pcpa.product_class_id = 2
	GROUP BY attr.name
) ch;

