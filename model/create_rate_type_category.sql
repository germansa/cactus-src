CREATE TABLE rate_type_categories
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    rate_type_uuid uuid NOT NULL,
    category_uuid uuid NOT NULL,
    amount FLOAT NOT NULL,
    percent FLOAT NOT NULL,
    mandatory BOOLEAN NOT NULL,
    active BOOLEAN NOT NULL,
    CONSTRAINT rate_type_categories_pkey PRIMARY KEY (id)
);
