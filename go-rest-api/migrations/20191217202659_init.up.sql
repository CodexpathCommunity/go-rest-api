CREATE TABLE users
(
    emailAddress         VARCHAR PRIMARY KEY,
    role       VARCHAR NOT NULL,
    name       VARCHAR NOT NULL,
    nickName       VARCHAR NOT NULL,
    country       VARCHAR NOT NULL,
    score       VARCHAR NOT NULL,
    createdAt TIMESTAMP NOT NULL,
    updatedAt TIMESTAMP NOT NULL
);