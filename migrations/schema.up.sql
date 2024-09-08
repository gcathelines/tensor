CREATE TABLE IF NOT EXISTS power_plants(
    "id" BIGSERIAL PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "latitude" NUMERIC NOT NULL,
    "longitude" NUMERIC NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NULL
);

