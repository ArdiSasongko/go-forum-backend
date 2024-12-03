-- +goose Up
-- +goose StatementBegin
alter table user_sessions
    alter COLUMN token type varchar(375);
-- +goose StatementEnd

-- +goose StatementBegin
alter table user_sessions
    alter COLUMN token set not null;
-- +goose StatementEnd

-- +goose StatementBegin
alter table user_sessions
    alter COLUMN refresh_token type varchar(375);
-- +goose StatementEnd

-- +goose StatementBegin
alter table user_sessions
    alter COLUMN refresh_token set not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table user_sessions
    alter COLUMN token type varchar(255);
-- +goose StatementEnd

-- +goose StatementBegin
alter table user_sessions
    alter COLUMN token set not null;
-- +goose StatementEnd

-- +goose StatementBegin
alter table user_sessions
    alter COLUMN refresh_token type varchar(255);
-- +goose StatementEnd

-- +goose StatementBegin
alter table user_sessions
    alter COLUMN refresh_token set not null;
-- +goose StatementEnd

