-- +goose Up
-- +goose StatementBegin
create table if not exists images_content (
    id serial not null,
    content_id int not null,
    image_url varchar(255) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    primary key (id),
    constraint fk_images_content_id foreign key (content_id) references contents(id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists images_content;
-- +goose StatementEnd
