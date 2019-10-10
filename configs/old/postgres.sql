
CREATE TABLE details (
  id SERIAL PRIMARY KEY,
  detail TEXT
);

INSERT INTO lapkin.details (detail)
VALUES ('авторы'), ('материал'), ('покрытие'), ('тип печати'), ('размер');

--*******************************************************

-- sizes
CREATE TABLE sizes (
  id SERIAL PRIMARY KEY,
  size TEXT
);

INSERT INTO lapkin.sizes (size)
VALUES ('105х148'), ('148x210'), ('130x180'), ('300x450'), ('300x450'), ('200x300'), ('400x600'), ('600x900'), ('1000x1500'), ('800x1200'), ('A4'), ('A6');

--*******************************************************

-- authors
CREATE TABLE authors (
  id SERIAL PRIMARY KEY,
  author TEXT
);

INSERT INTO lapkin.authors (author)
VALUES ('Анастасия Кондратьева'), ('Lolka Lolkina');

--*******************************************************

CREATE TABLE materials (
  id SERIAL PRIMARY KEY,
  material TEXT
);

INSERT INTO lapkin.materials (material)
VALUES ('240 g/m² pure white paper'), ('300 g/m² Munken Lynx Rough paper (woodfree)');

--*******************************************************

CREATE TABLE finishes (
  id SERIAL PRIMARY KEY,
  finish TEXT
);

INSERT INTO lapkin.finishes (finish)
VALUES ('semi-gloss'), ('gloss'), ('matte');

--*******************************************************

CREATE TABLE products_to_details (
  detail_id INT REFERENCES details(id),
	product_id INT REFERENCES products(id),
	value TEXT
);


CREATE TABLE print_types (
  id SERIAL PRIMARY KEY,
  print_type TEXT
);

INSERT INTO lapkin.print_types (print_type)
VALUES ('digital printing'), ('12-colour digital printing');

--*******************************************************

-- product_types
CREATE TABLE product_types (
  id SERIAL PRIMARY KEY,
  product_type TEXT
);

INSERT INTO lapkin.product_types (product_type)
VALUES ('postcards'), ('posters');

--*******************************************************

CREATE TABLE product_types_to_details (
	id SERIAL PRIMARY KEY,
  product_type_id INTEGER REFERENCES product_types(id),
	detail_id INTEGER REFERENCES details(id)
);

INSERT INTO lapkin.product_types_to_details (product_type_id, detail_id)
VALUES  (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (2, 1), (2, 2), (2, 3), (2, 4), (2, 5);

--*******************************************************

CREATE TABLE product_types_to_sizes (
	id SERIAL PRIMARY KEY,
  product_type_id INTEGER REFERENCES product_types(id),
	size_id INTEGER REFERENCES sizes(id)
);

INSERT INTO lapkin.product_types_to_sizes (product_type_id, size_id)
VALUES (2, 3), (2, 4), (2, 5), (2, 6), (2, 7), (2, 8);

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

CREATE TABLE products_to_authors (
	id SERIAL PRIMARY KEY,
  author_id INTEGER REFERENCES authors(id),
	product_id INTEGER REFERENCES products(id)
);

INSERT INTO lapkin.products_to_authors (author_id, product_id)
VALUES (1, 1), (2, 1), (2, 2), (1, 3), (1, 4), (2, 5), (1, 6), (2, 6);

--*******************************************************

CREATE TABLE products_to_sizes (
	id SERIAL PRIMARY KEY,
  size_id INTEGER REFERENCES sizes(id),
	product_id INTEGER REFERENCES products(id)
);

INSERT INTO lapkin.products_to_sizes (size_id, product_id)
VALUES (2, 1), (3, 1), (4, 1), (2, 2), (3, 2), (4, 3), (2, 4), (3, 4), (4, 4);

--*******************************************************

CREATE TABLE products_to_materials (
	id SERIAL PRIMARY KEY,
  material_id INTEGER REFERENCES materials(id),
	product_id INTEGER REFERENCES products(id)
);

INSERT INTO lapkin.products_to_materials (material_id, product_id)
VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6);

--*******************************************************F=

CREATE TABLE products_to_finishes (
	id SERIAL PRIMARY KEY,
  finish_id INTEGER REFERENCES finishes(id),
	product_id INTEGER REFERENCES products(id)
);

INSERT INTO lapkin.products_to_finishes (finish_id, product_id)
VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6);

--*******************************************************

CREATE TABLE products_to_print_types (
	id SERIAL PRIMARY KEY,
  print_type_id INTEGER REFERENCES print_types(id),
	product_id INTEGER REFERENCES products(id)
);

INSERT INTO lapkin.products_to_print_types (print_type_id, product_id)
VALUES (1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6);

--*******************************************************


with cte as (
	SELECT *
	FROM lapkin.products AS products
	WHERE product_type = 2
)
SELECT
	mmm.name, 
	mmm.description, 
	mmm.price, 
	mmm.is_available,
	json_build_object(
	    'authors', json_agg(mmm.author),
	    'finishes', json_agg(mmm.finish),
	    'materials', json_agg(mmm.material),
	    'print_types', json_agg(mmm.print_type)
	) AS details
from (
	SELECT
		x.name, 
		x.description, 
		x.price, 
		x.is_available,
	    a.author,
	    f.finish, 
	    m.material, 
  		pt.print_type
	FROM (
        select
            cte.*,
		    pta.author_id AS authors,
		    ptf.finish_id AS finishes, 
		    ptm.material_id AS materials, 
	  		ptpt.print_type_id AS print_types
		FROM cte
		JOIN lapkin.products_to_authors AS pta ON  cte.id = pta.product_id
		JOIN lapkin.products_to_finishes AS ptf ON  cte.id = ptf.product_id
		JOIN lapkin.products_to_materials AS ptm ON  cte.id = ptm.product_id
		JOIN lapkin.products_to_print_types AS ptpt ON  cte.id = ptpt.product_id
    ) x
    JOIN lapkin.authors AS a ON x.authors = a.id
	  JOIN lapkin.finishes AS f ON x.finishes = f.id
		JOIN lapkin.materials AS m ON x.materials = m.id
		JOIN lapkin.print_types AS pt ON x.print_types = pt.id
) mmm
GROUP BY mmm.name,
mmm.description,
mmm.price,
mmm.is_available;


-- ***************************************************


WITH cte AS (
	SELECT *
	FROM lapkin.products AS products
	WHERE product_type = 2
)
SELECT
	mmm.name,
	mmm.description, 
	mmm.price, 
	mmm.is_available,
	json_build_object(
	    mmm.detail, json_agg(mmm.value)
	) AS details
FROM (
	SELECT
		x.name,
		x.description, 
		x.price, 
		x.is_available,
		pdetails.value,
		d.detail
	FROM (
        SELECT
        	cte.id,
            cte.name,
			cte.description,
			cte.price, 
			cte.is_available,
		    ptd.detail_id AS detail
		FROM cte
		JOIN lapkin.product_types_to_details AS ptd ON cte.product_type = ptd.product_type_id
    ) x
    JOIN lapkin.products_to_details AS pdetails ON x.detail = pdetails.detail_id AND x.id = pdetails.product_id
    JOIN lapkin.details AS d ON x.detail = d.id
   
) mmm
GROUP BY mmm.name,
mmm.description,
mmm.price,
mmm.is_available,
mmm.detail;

-- ***************************************************

WITH cte AS (
	SELECT *
	FROM lapkin.products AS products
	WHERE product_type = 2
)
SELECT
	mmm.name,
	mmm.description, 
	mmm.price, 
	mmm.is_available,
	array_to_json(array_agg(row_to_json((select detail, value from mmm))))
--	json_build_object(
--	    mmm.detail, json_agg(mmm.value)
--	) AS details
FROM (
	SELECT
		x.name,
		x.description, 
		x.price, 
		x.is_available,
		x.value,
		x.detail
	FROM (
        SELECT
        	cte.id,
            cte.name,
			cte.description,
			cte.price, 
			cte.is_available,
			pdetails.value AS value,
			d.detail AS detail,
			array_to_json(array_agg(row_to_json((select detail, value from mmm))))
		FROM cte
		JOIN lapkin.product_types_to_details AS ptd ON cte.product_type = ptd.product_type_id
		LEFT OUTER JOIN lapkin.products_to_details AS pdetails ON ptd.detail_id = pdetails.detail_id AND pdetails.product_id = cte.id
   		JOIN lapkin.details AS d ON ptd.detail_id = d.id
    ) x
) mmm
GROUP BY mmm.name,
mmm.description,
mmm.price,
mmm.is_available;


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