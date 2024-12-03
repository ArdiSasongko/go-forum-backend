-- +goose Up
-- +goose StatementBegin
create table if not exists user_sessions (
    user_id int not null,
    token varchar(255) not null,
    token_expired timestamp not null,
    refresh_token varchar(255) not null,
    refresh_token_expired timestamp not null,
    constraint fk_user_id foreign key (user_id) references users(id) on delete cascade,
    constraint unique_user_id unique (user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_sessions;
-- +goose StatementEnd
