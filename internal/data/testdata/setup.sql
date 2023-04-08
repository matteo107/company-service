CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS company (
id uuid DEFAULT uuid_generate_v4(),
    name varchar(15) NOT NULL unique,
description varchar(3000) NULL,
employees integer NOT NULL,
registered boolean NOT NULL,
type text NOT NULL
);

INSERT INTO company (id,name, description, employees, registered, type) VALUES ('f1203d76-0491-47fe-9640-0aeda76ad3f6','Company One', 'Description for company one', 100, true, 'Corporations');