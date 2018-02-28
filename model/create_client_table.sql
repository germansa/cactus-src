CREATE TABLE clients
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    email VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
    CONSTRAINT clients_pkey PRIMARY KEY (id)
);
