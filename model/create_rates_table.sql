CREATE TABLE rates
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price FLOAT NOT NULL,
    rate_type_uuid uuid NOT NULL,
    reservation_uuid uuid NOT NULL,
    CONSTRAINT rates_pkey PRIMARY KEY (id)
);
