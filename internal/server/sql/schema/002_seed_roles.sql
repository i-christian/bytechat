-- +goose Up
INSERT INTO roles (name, description)
VALUES 
    ('admin', 'Full access to the system'),
    ('user', 'normal users who can send and receive messages')
ON CONFLICT (name) DO NOTHING;

-- +goose Down
DELETE FROM roles;
