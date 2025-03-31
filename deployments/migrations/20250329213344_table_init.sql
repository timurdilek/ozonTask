-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       author_id INT NOT NULL,
                       content TEXT NOT NULL,
                       are_comments_allowed BOOLEAN DEFAULT TRUE,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT '0001-01-01 00:00:00'
);

CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
                          parent_comment_id INT REFERENCES comments(id) ON DELETE CASCADE,
                          author_id INT NOT NULL,
                          content VARCHAR(2000) NOT NULL,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT '0001-01-01 00:00:00'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts CASCADE ;
DROP TABLE IF EXISTS comments;
-- +goose StatementEnd
