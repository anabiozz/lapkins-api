-- get_products
DROP FUNCTION IF EXISTS cart.add_product(text);
CREATE OR REPLACE FUNCTION cart.add_product(text)
RETURNS TABLE  (
	session_id uuid
)
AS $$
	BEGIN
	 	RETURN QUERY
	 	INSERT INTO 
 			cart.cart(data)
		VALUES
			($1::jsonb)
		RETURNING session;
	END;
$$ LANGUAGE plpgsql;