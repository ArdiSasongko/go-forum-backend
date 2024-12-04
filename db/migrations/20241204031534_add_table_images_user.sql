-- +goose Up
-- +goose StatementBegin
create table if not exists images_user (
    id serial not null,
    user_id int not null,
    image_url varchar(255) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    primary key (id),
    constraint fk_images_user_id foreign key (user_id) references users(id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists images_user;
-- +goose StatementEnd
