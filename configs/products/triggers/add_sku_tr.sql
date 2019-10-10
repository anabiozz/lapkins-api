DROP TRIGGER IF EXISTS "add_sku_tr" ON products.product_variant;

CREATE OR REPLACE FUNCTION products.add_sku()
RETURNS TRIGGER AS $$
BEGIN
	NEW.sku = (
    	SELECT pc.sku_part || '-' || currval(pg_get_serial_sequence(TG_TABLE_NAME, 'id'))
	   	FROM products.product_class AS pc
	   	JOIN products.product AS p ON p.product_class_id = pc.id
	   	JOIN (SELECT NEW.*) AS pv ON pv.product_id = p.id AND pv.id = currval(pg_get_serial_sequence(TG_TABLE_NAME, 'id'))
	);
    RETURN NEW;   
END;
$$ language 'plpgsql';

CREATE TRIGGER add_sku_tr BEFORE INSERT ON products.product_variant FOR EACH ROW EXECUTE PROCEDURE products.add_sku();