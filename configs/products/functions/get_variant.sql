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