-- +goose Up
-- +goose StatementBegin
create type Roles as enum ('admin', 'user');
-- +goose StatementEnd

-- +goose StatementBegin
create table if not exists users (
    id serial not null,
    name varchar(255) not null,
    username varchar(255) not null unique,
    email varchar(255) not null,
    password varchar(255) not null,
    role Roles not null,
    is_valid boolean,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    primary key (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd

-- +goose StatementBegin
drop type if exists Roles;
-- +goose StatementEnd