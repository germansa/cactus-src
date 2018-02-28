CREATE TABLE histories
(
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    table_reference VARCHAR NOT NULL,
    field VARCHAR NOT NULL,
    value VARCHAR NOT NULL,
    created_date VARCHAR NOT NULL,
    user_uuid uuid NOT NULL,
    row_uuid uuid NOT NULL,
    CONSTRAINT histories_pkey PRIMARY KEY (id)
);
