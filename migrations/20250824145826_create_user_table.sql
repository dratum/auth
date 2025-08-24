-- +goose Up
create table users (
    id serial primary key,
    name text not null,
    email text not null unique,
    password_hash text not null,
    role text not null,
    created_at timestamp not null default now(),
    updated_at TIMESTAMP
);
-- +goose Down
drop table user;