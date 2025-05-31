-- +goose Up
-- +goose StatementBegin
ALTER TABLE books 
    ADD COLUMN id Serial PRIMARY KEY,
    ADD COLUMN title VARCHAR(255) NOT NULL,
    ADD COLUMN author VARCHAR(255) NOT NULL,
    ADD COLUMN summary TEXT,
    ADD COLUMN created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE books;
-- +goose StatementEnd