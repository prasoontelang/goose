-- +goose Up
CREATE TABLE users (
                       id int NOT NULL PRIMARY KEY,
                       username text,
                       name text,
                       surname text
);

-- +goose Down
DROP TABLE users;