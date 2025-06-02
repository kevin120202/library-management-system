-- +goose Up
-- +goose StatementBegin
ALTER TABLE books 
    ADD COLUMN availability BOOLEAN
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE books
DROP COLUMN availability
-- +goose StatementEnd