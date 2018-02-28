CREATE TABLE rate_types
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    priority INTEGER NOT NULL,
    CONSTRAINT rate_types_pkey PRIMARY KEY (id)
);
