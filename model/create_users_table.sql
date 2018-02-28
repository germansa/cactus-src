CREATE TABLE users
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
    email VARCHAR NOT NULL UNIQUE,
    CONSTRAINT userss_pkey PRIMARY KEY (id)
);
