-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_activities (
    id SERIAL NOT NULL,
    user_id INT NOT NULL,
    content_id INT NOT NULL,
    isLiked BOOLEAN DEFAULT FALSE,
    isDisliked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    created_by VARCHAR(255) NOT NULL,
    updated_by VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_user_activities_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_activities_content_id FOREIGN KEY (content_id) REFERENCES contents(id) ON DELETE CASCADE,
    CONSTRAINT chk_user_activities_check CHECK (
        (isLiked = TRUE AND isDisliked = FALSE) OR
        (isLiked = FALSE AND isDisliked = TRUE) OR
        (isLiked = FALSE AND isDisliked = FALSE)
    ),
    CONSTRAINT user_content_unique UNIQUE (user_id, content_id) -- Tambahkan constraint unik
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_activities;
-- +goose StatementEnd
