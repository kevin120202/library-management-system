-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS books (
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books;
-- +goose StatementEnd