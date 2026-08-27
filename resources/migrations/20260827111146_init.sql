-- +goose Up
CREATE TABLE "products" (
    "id" UUID PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE "product_images" (
    "id" UUID PRIMARY KEY,
    "path" TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS "product_images";
DROP TABLE IF EXISTS "products";
