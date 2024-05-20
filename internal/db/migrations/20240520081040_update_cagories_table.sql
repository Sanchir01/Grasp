-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories
    ADD COLUMN slug VARCHAR(255) UNIQUE,
    ADD COLUMN description VARCHAR(255) ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
