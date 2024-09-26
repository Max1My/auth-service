-- +goose Up
CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    role_id       BIGINT,
    telegram_id   BIGINT UNIQUE,
    email VARCHAR(255) UNIQUE,
    email_verified BOOLEAN DEFAULT false,
    FOREIGN KEY (role_id) REFERENCES roles (id),
    CONSTRAINT chk_password_or_telegram_id CHECK (
        (password_hash IS NOT NULL AND telegram_id IS NULL) OR
        (telegram_id IS NOT NULL AND password_hash IS NULL)
        )
);

-- +goose Down
DROP TABLE users;
