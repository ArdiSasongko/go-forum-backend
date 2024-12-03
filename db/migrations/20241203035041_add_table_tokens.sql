-- +goose Up
-- +goose StatementBegin
create type token_type as enum ('email', 'password_reset');
-- +goose StatementEnd

-- +goose StatementBegin
create table if not exists tokens (
    user_id int not null,
    token_type token_type not null,
    token int not null,
    expired_at timestamp not null,
    constraint fk_token_user foreign key (user_id) references users(id) on delete cascade,
    constraint unique_token_user_id unique (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists tokens;
-- +goose StatementEnd

-- +goose StatementBegin
drop type if exists token_type;
-- +goose StatementEnd
