DROP FUNCTION IF EXISTS cart.increase_quantity(text, text, text);
CREATE OR REPLACE FUNCTION cart.increase_quantity(text, text, text)
RETURNS TABLE  (
	session_id uuid
)
AS $$
	BEGIN
	 	RETURN QUERY
	 	UPDATE 
	 		cart.cart
		SET
			"data" = "data"||'{"quantity"}'
		WHERE
	END;
$$ LANGUAGE plpgsql;
