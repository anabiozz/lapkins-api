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