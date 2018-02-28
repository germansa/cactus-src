CREATE TABLE categories
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    image VARCHAR NOT NULL,
    CONSTRAINT categories_pkey PRIMARY KEY (id)
);
