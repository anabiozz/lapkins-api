DROP FUNCTION IF EXISTS cart.create_session();
CREATE OR REPLACE FUNCTION cart.create_session()
RETURNS TABLE  ( session_id uuid )
AS $$
	BEGIN
	RETURN QUERY
		INSERT INTO
 			cart.cart("session")
		VALUES
			((SELECT uuid_generate_v4()))
		RETURNING session;
	END;
$$ LANGUAGE plpgsql;

SELECT * FROM cart.create_session()