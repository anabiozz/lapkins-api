DROP TRIGGER IF EXISTS "update_modtime" ON products.product;
DROP TRIGGER IF EXISTS "update_modtime" ON products.product_variant;

CREATE OR REPLACE FUNCTION products.update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;   
END;
$$ language 'plpgsql';

CREATE TRIGGER update_modtime BEFORE UPDATE ON products.product FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER update_modtime BEFORE UPDATE ON products.product_variant FOR EACH ROW EXECUTE PROCEDURE update_modified_column();