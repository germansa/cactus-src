CREATE TABLE reservations
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    confirmation_code VARCHAR NOT NULL,
    client_uuid uuid NOT NULL,
    status INTEGER NOT NULL,
    pickup_model_uuid uuid NOT NULL,
    category_uuid uuid NOT NULL,
    created_date VARCHAR NOT NULL,
    total FLOAT,
    pickup_date VARCHAR NOT NULL,
    pickup_time VARCHAR NOT NULL,
    return_date VARCHAR NOT NULL,
    return_time VARCHAR NOT NULL,
    CONSTRAINT reservations_pkey PRIMARY KEY (id)
);
