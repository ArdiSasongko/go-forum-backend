-- +goose Up
-- +goose StatementBegin
create table if not exists comments (
    id serial not null,
    user_id int not null,
    content_id int not null,
    comment_body text not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    created_by varchar(255) not null,
    updated_by varchar(255) not null,
    primary key (id),
    constraint fk_comments_user_id foreign key (user_id) references users(id) on delete cascade,
    constraint fk_comments_content_id foreign key (content_id) references contents(id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists comments;
-- +goose StatementEnd
