-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          price INT NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
