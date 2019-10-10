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