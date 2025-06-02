-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS borrows_returns (
    id BIGSERIAL PRIMARY KEY,
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE borrows_returns;
-- +goose StatementEnd