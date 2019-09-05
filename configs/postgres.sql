--CREATE SCHEMA lapkin;

--ALTER TABLE lapkin.details_to_products DROP COLUMN detail;

--*******************************************************

-- sizes
CREATE TABLE sizes (
  id SERIAL PRIMARY KEY,
  value TEXT
);

INSERT INTO lapkin.sizes (value)
VALUES ('105х148'), ('148x210'), ('130x180'), ('300x450'), ('300x450'), ('200x300'), ('400x600'), ('600x900'), ('1000x1500'), ('800x1200');

--*******************************************************

-- authors
CREATE TABLE authors (
  id SERIAL PRIMARY KEY,
  author TEXT
);

INSERT INTO lapkin.authors (author)
VALUES ('Анастасия Кондратьева'), ('Lolka Lolkina');

--*******************************************************

CREATE TABLE details (
  id SERIAL PRIMARY KEY,
  details TEXT
);

INSERT INTO lapkin.details (details)
VALUES ('авторы'), ('материал'), ('покрытие'), ('тип печати'), ('размер');

--*******************************************************

-- product_types
CREATE TABLE product_types (
  id SERIAL PRIMARY KEY,
  product_type TEXT
);

INSERT INTO lapkin.product_types (product_type)
VALUES ('postcards'), ('posters');

--*******************************************************

CREATE TABLE details_to_product_types (
	id SERIAL PRIMARY KEY,
  	product_type_id INTEGER REFERENCES product_types(id),
	detail_id INTEGER REFERENCES details(id)
);

INSERT INTO lapkin.details_to_product_types (product_type_id, detail_id)
VALUES (2, 1), (2, 2), (2, 3), (2, 4), (2, 5);

--*******************************************************

-- products
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name TEXT,
  description TEXT,
  price INTEGER,
  product_type INTEGER REFERENCES product_types(id),
  is_available BOOLEAN
);

INSERT INTO lapkin.products (name, description, price, product_type, is_available)
VALUES ('product_1', 'descripttion 1', 50, 2, true), 
			 ('product_2', 'descripttion 2', 50, 2, true), 
			 ('product_3', 'descripttion 3', 50, 2, true), 
			 ('product_4', 'descripttion 4', 50, 2, true),
			 ('product_5', 'descripttion 5', 50, 2, true), 
			 ('product_6', 'descripttion 6', 50, 2, true);
			
--*******************************************************

CREATE TABLE authors_to_products (
	id SERIAL PRIMARY KEY,
  	author_id INTEGER REFERENCES authors(id),
	product_id INTEGER REFERENCES products(id)
);

INSERT INTO lapkin.authors_to_products (author_id, product_id)
VALUES (1, 1), (2, 1);

--*******************************************************

CREATE TABLE sizes_to_products (
	id SERIAL PRIMARY KEY,
  	size_id INTEGER REFERENCES sizes(id),
	product_id INTEGER REFERENCES products(id)
);

INSERT INTO lapkin.sizes_to_products (size_id, product_id)
VALUES (2, 1), (3, 1), (4, 1);

--*******************************************************

CREATE TABLE details_to_products (
	id SERIAL PRIMARY KEY,
  	details_id INTEGER REFERENCES details(id),
	product_id INTEGER REFERENCES products(id)
);

INSERT INTO lapkin.details_to_products (details_id, product_id)
VALUES (1, 1), (2, 1), (3, 1), (4, 1), (5, 1);

--*******************************************************


SELECT name, description, price, is_available, dtp.details_id
FROM lapkin.products AS p
JOIN lapkin.details_to_products AS dtp ON  p.id = dtp.product_id 
WHERE product_type = 2


-- get_products
CREATE OR REPLACE FUNCTION lapkin.get_products(INT)
RETURNS SETOF lapkin.products
AS $$
	BEGIN
	 	RETURN QUERY SELECT * FROM lapkin.products WHERE product_type = $1;
	END;
$$ LANGUAGE plpgsql;

-- get_authors
CREATE OR REPLACE FUNCTION lapkin.get_authors()
RETURNS SETOF lapkin.authors
AS $$
	BEGIN
	 	RETURN QUERY SELECT * FROM lapkin.authors;
	END;
$$ LANGUAGE plpgsql;

-- get_product_types
CREATE OR REPLACE FUNCTION lapkin.get_product_types()
RETURNS SETOF lapkin.product_types
AS $$
	BEGIN
	 	RETURN QUERY SELECT * FROM lapkin.product_types;
	END;
$$ LANGUAGE plpgsql;

-- get_sizes
CREATE OR REPLACE FUNCTION lapkin.get_sizes(INT)
RETURNS TABLE  (
	size_id int,
	proportions text
)
AS $func$
	BEGIN
		RETURN QUERY 
 		SELECT sizes.size_id, sizes.proportions FROM (
	      SELECT sizes.size_id, sizes.proportions 
	      FROM   lapkin.sizes
	      WHERE product_type = $1
     	) sizes;
	END;
$func$ LANGUAGE plpgsql;

-- get_product_by_id
CREATE OR REPLACE FUNCTION lapkin.get_product_by_id(INT)
RETURNS SETOF lapkin.products
AS $$
	BEGIN
	 	RETURN QUERY SELECT * FROM lapkin.products WHERE id = $1;
	END;
$$ LANGUAGE plpgsql;