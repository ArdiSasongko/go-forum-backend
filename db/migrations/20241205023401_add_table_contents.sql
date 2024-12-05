-- +goose Up
-- +goose StatementBegin
create table if not exists contents (
    id serial not null,
    user_id int not null,
    content_title varchar(255) not null,
    content_body text not null,
    content_hastags varchar(255) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    created_by varchar(255) not null,
    updated_by varchar(255) not null,
    primary key (id),
    constraint fk_content_user_id foreign key (user_id) references users(id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists contents;
-- +goose StatementEnd
