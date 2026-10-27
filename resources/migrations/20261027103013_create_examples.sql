-- +goose Up
CREATE TABLE IF NOT EXISTS "examples" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "title" TEXT NOT NULL,
    "done" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS "examples_open_idx" ON "examples" ("created_at" DESC) WHERE NOT "done";

-- +goose Down
DROP TABLE IF EXISTS "examples";
