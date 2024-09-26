-- +goose Up
CREATE TABLE roles
(
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE roles;
