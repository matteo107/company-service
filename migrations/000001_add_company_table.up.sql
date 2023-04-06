CREATE TABLE IF NOT EXISTS company (
id uuid DEFAULT uuid_generate_v4(),
name varchar(15) NOT NULL unique,
description varchar(3000) NULL,
employees integer NOT NULL,
registered boolean NOT NULL,
type text NOT NULL
);