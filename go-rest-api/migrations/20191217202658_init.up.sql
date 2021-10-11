CREATE TABLE idea
(
    id                  VARCHAR PRIMARY KEY,
    author_email        VARCHAR NOT NULL,
    tags                text[],
    summary             VARCHAR,
    media               text[],
    bad_flag             boolean,
    enabled             boolean,
    issues              text[],
    votes               integer,
    created_at          TIMESTAMP NOT NULL,
    updated_at          TIMESTAMP NOT NULL
);

