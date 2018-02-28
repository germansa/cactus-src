CREATE TABLE models
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    make VARCHAR NOT NULL,
    category_uuid uuid NOT NULL,
    CONSTRAINT models_pkey PRIMARY KEY (id)
);
