CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS lords
(
    id              uuid            PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            varchar(200)    NOT NULL,
    tv_series       varchar[]       NOT NULL,
    created_at      TIMESTAMP       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS lords_pkey ON lords USING btree (id);
CREATE UNIQUE INDEX IF NOT EXISTS lords_name ON lords USING btree (name,deleted_at);

CREATE TABLE IF NOT EXISTS houses
(
    id                  uuid            PRIMARY KEY DEFAULT uuid_generate_v4(),
    name                varchar(200)    NOT NULL    UNIQUE,
    region              varchar(100)    NOT NULL,
    foundation_year     varchar(5)      NOT NULL,
    current_lord        uuid,
    created_at          TIMESTAMP       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP       NOT NULL    DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP,
    CONSTRAINT fk_lords
        FOREIGN KEY(current_lord)
            REFERENCES lords(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS houses_pkey ON houses USING btree (id);
CREATE UNIQUE INDEX IF NOT EXISTS houses_name ON houses USING btree (name,deleted_at);