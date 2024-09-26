-- +goose Up
CREATE TABLE password_reset_tokens
(
    id         SERIAL PRIMARY KEY,
    user_id    BIGINT,
    token      VARCHAR(255) NOT NULL,
    created_at timestamp    not null default now(),
    expires_at TIMESTAMP    NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
-- +goose Down
DROP TABLE password_reset_tokens;
