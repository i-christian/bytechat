-- +goose Up
insert into roles (name, description)
values
    ('admin', 'Full access to the system'),
    ('user', 'normal users who can send and receive messages')
on conflict (name) do nothing;

-- +goose Down
delete from roles;
