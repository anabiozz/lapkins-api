CREATE SCHEMA IF NOT EXISTS products AUTHORIZATION lapkin;

-- FUNCTIONS ********************************************

CREATE OR REPLACE FUNCTION products.update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.product_id = now();
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
VALUES ('authors', 'авторы'), ('materials', 'материалы'), ('finish', 'покрытие'), ('print type', 'тип печати'), ('size', 'размер');


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
	attribute_id INT REFERENCES products.product_attribute(id),
	product_class_id INT REFERENCES products.product_class(id)
);

INSERT INTO products.attribute_choice_value (display, attribute_id, product_class_id)
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
VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (2, 1), (2, 2), (2, 3), (2, 4), (2, 5);

CREATE TABLE products.product_class_variant_attribute (
	id SERIAL PRIMARY KEY,
	product_class_id INT REFERENCES products.product_class(id),
	product_attribute_id INT REFERENCES products.product_attribute(id)
);

INSERT INTO products.product_class_variant_attribute (product_class_id, product_attribute_id)
VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (2, 1), (2, 2), (2, 3), (2, 4), (2, 5);

-- PRODUCT ********************************************

CREATE TABLE products.product (
	id SERIAL PRIMARY KEY,
	name TEXT,
	description TEXT,
	price INT,
	weight INT,
	available_on BOOLEAN,
	updated_at timestamptz DEFAULT current_timestamp,
	product_class_id INT REFERENCES product_class(id)
--	attributes JSONB
);

INSERT INTO products.product (name, description, price, product_class_id)
VALUES ('плакат веселье', 'Et deserunt labore excepteur id eiusmod reprehenderit do nostrud cupidatat consectetur laboris culpa.', 300, 2),
('плакат надпись со смыслом', 'Commodo labore est qui laboris irure esse aliquip', 300, 2),
('открытка веселье', 'Anim ex occaecat occaecat non tempor in enim id mollit.', 50, 1),
('открытка надпись со смыслом', 'Nisi eiusmod laborum ullamco mollit elit amet deserunt ex sit nisi consectetur cillum commodo incididunt.', 300, 1);


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

INSERT INTO products.product_variant (name, price_override, product_id)
VALUES ('плакат веселье 300x450', 300, 1), 
('плакат веселье с металической рамой 300x450', 600, 1),
('плакат веселье с деревянной рамой 300x450', 400, 1),
('плакат веселье 600x900', 300, 1), 
('плакат веселье с металической рамой 600x900', 600, 1),
('плакат веселье с деревянной рамой 600x900', 400, 1),
('открытка веселье 105х148', 50, 3),
('плакат надпись со смыслом 300x450', 300, 2), 
('плакат надпись со смыслом с металической рамой 300x450', 600, 2),
('плакат надпись со смыслом с деревянной рамой 300x450', 400, 2),
('плакат надпись со смыслом 600x900', 300, 2), 
('плакат надпись со смыслом с металической рамой 600x900', 600, 2),
('плакат надпись со смыслом с деревянной рамой 600x900', 400, 2),
('открытка надпись со смыслом 105х148', 50, 4);


CREATE TABLE products.product_attribute_choice (
	id SERIAL PRIMARY KEY,
	product_id INT REFERENCES products.product(id),
	attribute_choice_id INT REFERENCES products.attribute_choice_value(id)
)


CREATE TABLE products.product_variant_attribute (
	id SERIAL PRIMARY KEY,
	product_varian_id INT REFERENCES products.product_variant(id),
	attribute_choice_id INT REFERENCES products.attribute_choice_value(id)
);

INSERT INTO products.product_variant_attribute (product_varian_id, attribute_choice_id)
VALUES (1, 13), (1, 1), (1, 2), (1, 8), (1, 5), (1, 3), (7, 10), (7, 1), (7, 2), (7, 8), (7, 5), (7, 3);



UPDATE products.product_variant
SET ATTRIBUTES = (
WITH attr AS (
	SELECT * FROM product_attribute
)
SELECT json_object_agg(ch.name, ch.display)
FROM (
	SELECT attr.name, json_agg(acv.display) AS display
	FROM attr
	JOIN attribute_choice_value AS acv ON attr.id = acv.attribute_id
	JOIN product_class_variant_attribute AS pcva ON attr.id = pcva.product_attribute_id AND pcva.product_class_id = 2
	GROUP BY attr.name
) ch)
WHERE id = 1;

CREATE TABLE products.product_categories (
	id SERIAL PRIMARY KEY,
	product_id INT REFERENCES products.product(id),
	category_id INT REFERENCES products.category(id)
);

INSERT INTO products.product_categories (product_id, category_id)
VALUES (1, 1), (2, 1), (3, 1), (4, 1);

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

SELECT * FROM product_class_variant_attribute
JOIN attribute_choice_value AS attr ON attr.attribute_id = product_class_variant_attribute.product_attribute_id;


WITH attr AS (
	SELECT * FROM product_attribute
)
SELECT json_object_agg(ch.name, ch.display)
FROM (
	SELECT attr.name, json_agg(acv.display) AS display
	FROM attr
	JOIN attribute_choice_value AS acv ON attr.id = acv.attribute_id
	JOIN product_class_variant_attribute AS pcva ON attr.id = pcva.product_attribute_id AND pcva.product_class_id = 2 
	GROUP BY attr.name
) ch;


WITH tmp AS (
    SELECT 
        option_type,
        json_agg(sto.name) as training_options
    FROM 
        safety_training_options as sto
    GROUP BY 
        sto.option_type
)
SELECT json_object_agg(option_type, training_options) FROM tmp

SELECT 
  a.id, a.name, 
  ( 
    SELECT json_agg(item)
    FROM (
      SELECT b.c1 AS x, b.c2 AS y 
      FROM b WHERE b.item_id = a.id
    ) item
  ) AS items
FROM a;

with t as (
  select
    title,
    body,
    published_at,
    'https://til.hashrocket.com/posts/' || slug as permalink,
    channels.name,
    developers.username
  from posts
  join channels on channels.id = posts.channel_id
  join developers on developers.id = posts.developer_id
)
select json_agg(t) from t;


SELECT to_json(sub) AS container_with_things
FROM  (
   SELECT c.*, json_agg(thing_id) AS things
   FROM   attr c
   LEFT   JOIN container_thing ct ON  ct.container_id = c.id
   WHERE  c.id IN (<list of container ids>)
   GROUP  BY c.id
   ) sub;
  

SELECT p."name", pa.display, acv.display FROM product_attribute AS pa
--JOIN product_class_variant_attribute AS pcva ON pcva.product_attribute_id = pa.id
JOIN product_class_product_attribute AS pcpa ON pcpa.product_attribute_id = pa.id
JOIN product_class AS pc ON pc.id = pcpa.product_class_id
JOIN product AS p ON p.product_class_id = pc.id
JOIN product_variant AS pv ON pv.product_id = p.id
JOIN attribute_choice_value AS acv ON acv.attribute_id = pa.id
WHERE pv.id = 1


SELECT pa.display, acv.display FROM product_variant AS pv
JOIN products.product AS p ON p.id = pv.product_id
JOIN product_class AS pc ON pc.id = p.product_class_id
JOIN product_class_variant_attribute AS pcva ON pcva.product_attribute_id = 

SELECT * FROM product_variant AS pv
JOIN products.product AS p ON p.id = pv.product_id
JOIN 

