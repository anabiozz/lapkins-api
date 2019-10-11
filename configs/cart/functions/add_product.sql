-- get_products
DROP FUNCTION IF EXISTS cart.add_product(int, int, int);
CREATE OR REPLACE FUNCTION cart.add_product(variant_id_p int, session_id_p int, customer_id_p int)
RETURNS TABLE  ( session_id uuid )
AS $$
	BEGIN
		IF customer_id_p = 0 THEN
			RETURN QUERY
			IF EXISTS(SELECT 1 FROM cart.cart WHERE "session" = )
			INSERT INTO
	 			cart.cart("session", variant_id, session_id_p)
			VALUES
				(uuid_generate_v4(), variant_id_p, quantit_p)
			RETURNING session;
		ELSE
			INSERT INTO
	 			cart.cart(customer_id, variant_id, quantity)
			VALUES
				(customer_id_p, variant_id_p, quantit_p);
		END IF;
	END;
$$ LANGUAGE plpgsql;