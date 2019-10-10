
SELECT *
FROM product_variant
WHERE '"600x900"' IN (SELECT jsonb_array_elements(product_variant."attributes" -> 'sizes')::jsonb::text)

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
	'{600x900}' IN 
		(
			SELECT ARRAY_AGG(s.x || 'x' || s.y) AS size
			FROM product_class_variant_size AS pcvs
			JOIN products."size" AS s ON s.id = pcvs.product_size_id
			WHERE pcvs.variant_id = 4
			
		)
AND pv."attributes"->'frame' IS NULL 
AND pv.product_id = 1
GROUP BY pv.id, p.id;


SELECT 
	product_variant.id,
	p.id,
	product_variant."name",
	p.description,
	product_variant.price_override, 
	product_variant."attributes",
	ARRAY_AGG(s.x || 'x' || s.y) AS size
FROM products.product_variant AS pv
JOIN products.product AS p ON pv.product_id = p.id
JOIN products.product_class_product_size AS pcps ON pcps.product_class_id = p.product_class_id
JOIN products."size" AS s ON s.id = pcps.product_size_id
WHERE  '' || p_size || '' 
	IN 
		(
			SELECT ARRAY_AGG(s.x || 'x' || s.y) AS size
			FROM products."size"
			JOIN products.product_class_variant_size AS pcvs ON pcvs.variant_id = pv.id
--					SELECT jsonb_array_elements(product_variant."attributes" -> 'sizes')::jsonb::text
		) 
AND product_variant."attributes"->'frame' IS NULL 
AND product_variant.product_id = p_id
GROUP BY product_variant.id, p.id;

-- **************************************************************************


SELECT * 
FROM product_variant 
WHERE "attributes" @> '{"sizes": ["600x900"]}' AND "attributes"->'frame' IS NULL AND product_id = 2

UPDATE products.product
SET ATTRIBUTES = (
WITH attr AS (
	SELECT pcva.product_attribute_id AS id
	FROM product_class_variant_attribute AS pcva
	WHERE pcva.product_class_id = 2
	UNION
	SELECT pcpa.product_attribute_id AS id
	FROM product_class_product_attribute AS pcpa
	WHERE pcpa.product_class_id = 2
)
SELECT json_object_agg(ch.name, ch.display)
FROM (
	SELECT p_attr.name, json_agg(acv.display) AS display
	FROM attr
	JOIN product_attribute AS p_attr ON p_attr.id = attr.id
	JOIN attribute_choice_value AS acv ON acv.attribute_id = attr.id
	GROUP BY p_attr.name
) ch)
WHERE product_class_id = 2;


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


WITH inserted AS (
  INSERT INTO CATEGORY (NAME) VALUES ('PARENT_1')
  RETURNING id
) INSERT INTO CATEGORY (NAME, PARENT_ID) VALUES
  ('CHILD_1_1', (SELECT inserted.id FROM inserted)),
  ('CHILD_2_1', (SELECT inserted.id FROM inserted))