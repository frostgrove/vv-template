-- +goose Up
-- The database owns the keys and the timestamps: the models tag them
-- `auto` / `generated`, so vv never writes them and reads them back.
CREATE TABLE IF NOT EXISTS "products" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "name" TEXT NOT NULL,
    "description" TEXT,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "product_images" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "path" TEXT NOT NULL
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION "touch_updated_at"() RETURNS TRIGGER AS $$
BEGIN
    NEW."updated_at" = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER "products_touch_updated_at"
    BEFORE UPDATE ON "products"
    FOR EACH ROW EXECUTE FUNCTION "touch_updated_at"();

-- +goose Down
DROP TRIGGER IF EXISTS "products_touch_updated_at" ON "products";
DROP FUNCTION IF EXISTS "touch_updated_at"();
DROP TABLE IF EXISTS "product_images";
DROP TABLE IF EXISTS "products";
