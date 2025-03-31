package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"ozon/internal/transport/graph/model"
	"ozon/pkg/logger"
	"ozon/pkg/postgresql"
	"time"
)

func NewPsql(ctx context.Context, DBConn string) *PsqlPool {

	log := logger.GetLogger()
	pool := postgresql.ConnectionDB(ctx, DBConn)

	return &PsqlPool{Pool: pool, Logger: log}
}

type PsqlPool struct {
	*pgxpool.Pool
	*zap.Logger
}

type times struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t times) parseTime() (string, string) {
	return t.CreatedAt.Format(time.DateTime), t.UpdatedAt.Format(time.DateTime)
}

func (p PsqlPool) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {

	var output = model.Post{
		AuthorID:           input.AuthorID,
		Content:            input.Content,
		AreCommentsAllowed: input.AreCommentsAllowed,
	}

	t := times{}

	query := "INSERT INTO posts (author_id, content, are_comments_allowed) VALUES ($1, $2, $3) RETURNING id, created_at"

	err := p.Pool.QueryRow(ctx, query, input.AuthorID, input.Content, input.AreCommentsAllowed).Scan(&output.ID, &t.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("PsqlPool insert post: %w", err)
	}

	output.CreatedAt, output.UpdatedAt = t.parseTime()
	return &output, nil

}

func (p PsqlPool) PostComment(ctx context.Context, input model.PostCommentInput) (*model.Comment, error) {

	var output = model.Comment{
		PostID:          input.PostID,
		ParentCommentID: input.ParentCommentID,
		AuthorID:        input.AuthorID,
		Content:         input.Content,
	}

	t := times{}

	query := ""
	var err error

	if input.ParentCommentID != nil {
		query = "INSERT INTO comments (post_id, parent_comment_id, author_id,content) VALUES ($1, $2, $3,$4) RETURNING id, created_at"

		err = p.Pool.QueryRow(ctx, query, input.PostID, input.ParentCommentID, input.AuthorID, input.Content).Scan(&output.ID, &t.CreatedAt)
	} else {
		query = "INSERT INTO comments (post_id, author_id,content) VALUES ($1, $2, $3) RETURNING id, created_at"

		err = p.Pool.QueryRow(ctx, query, input.PostID, input.AuthorID, input.Content).Scan(&output.ID, &t.CreatedAt)
	}

	if err != nil {
		return nil, fmt.Errorf("PsqlPool insert comments: %w", err)
	}
	output.CreatedAt, output.UpdatedAt = t.parseTime()

	return &output, err
}

func (p PsqlPool) PutPost(ctx context.Context, input model.PutPostInput) (*model.Post, error) {

	var output = model.Post{}

	t := times{}

	query := "UPDATE posts SET content = COALESCE($1, content), are_comments_allowed = COALESCE($2, are_comments_allowed), updated_at = NOW() WHERE id = $3 RETURNING *"

	err := p.Pool.QueryRow(ctx, query, input.Content, input.AreCommentsAllowed, input.ID).Scan(&output.ID, &output.AuthorID, &output.Content, &output.AreCommentsAllowed, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("PsqlPool update posts %w", err)
	}

	output.CreatedAt, output.UpdatedAt = t.parseTime()

	return &output, nil
}

func (p PsqlPool) PutComment(ctx context.Context, input model.PutCommentInput) (*model.Comment, error) {

	var output = model.Comment{
		ID:      input.ID,
		Content: input.Content,
	}

	t := times{}

	query := "UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2 RETURNING post_id,parent_comment_id,author_id,created_at,updated_at"

	err := p.Pool.QueryRow(ctx, query, input.Content, input.ID).Scan(&output.PostID, &output.ParentCommentID, &output.AuthorID, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("PsqlPool update comments %w", err)
	}
	output.CreatedAt, output.UpdatedAt = t.parseTime()

	return &output, err
}

func (p PsqlPool) DeletePost(ctx context.Context, id string) (bool, error) {

	query := "DELETE FROM posts WHERE id = $1"

	_, err := p.Pool.Exec(ctx, query, id)

	switch {
	case errors.Is(err, nil):
	case errors.Is(err, pgx.ErrNoRows):
		return false, nil
	default:
		return false, fmt.Errorf("PsqlPool delete post %w", err)
	}
	return true, err
}

func (p PsqlPool) DeleteComment(ctx context.Context, id string) (bool, error) {

	query := "DELETE FROM comments WHERE id = $1"

	_, err := p.Pool.Exec(ctx, query, id)

	switch {
	case errors.Is(err, nil):
	case errors.Is(err, pgx.ErrNoRows):
		return false, nil
	default:
		return false, fmt.Errorf("PsqlPool delete comments %w", err)
	}
	return true, err
}

func (p PsqlPool) GetPost(ctx context.Context, first int32) ([]*model.Post, error) {
	t := times{}

	var output []*model.Post

	query := "SELECT * FROM posts ORDER BY created_at DESC LIMIT 10 OFFSET $1"

	rows, err := p.Pool.Query(ctx, query, first)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("PsqlPool select posts %w", err)
	}

	var row model.Post

	for rows.Next() {

		err = rows.Scan(&row.ID, &row.AuthorID, &row.Content, &row.AreCommentsAllowed, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("PsqlPool select posts %w", err)
		}
		row.CreatedAt, row.UpdatedAt = t.parseTime()

		output = append(output, &row)
	}

	return output, nil
}

func (p PsqlPool) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	t := times{}

	var output model.Post

	query := "SELECT * FROM posts WHERE id = $1"

	err := p.Pool.QueryRow(ctx, query, id).Scan(&output.ID, &output.AuthorID, &output.Content, &output.AreCommentsAllowed, &t.CreatedAt, &t.UpdatedAt)
	switch {
	case errors.Is(err, nil):
	case errors.Is(err, pgx.ErrNoRows):
		return nil, nil
	default:
		return nil, fmt.Errorf("PsqlPool select post %w", err)
	}
	output.CreatedAt, output.UpdatedAt = t.parseTime()

	return &output, nil
}

func (p PsqlPool) GetCommentByPostID(ctx context.Context, postID string, first int32) ([]*model.Comment, error) {
	t := times{}

	var output []*model.Comment

	query := "SELECT * FROM comments WHERE post_id = $1 ORDER BY created_at DESC LIMIT 10 OFFSET $2"

	rows, err := p.Pool.Query(ctx, query, postID, first)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("PsqlPool select comment by postID %w", err)
	}

	var row model.Comment

	for rows.Next() {

		err = rows.Scan(&row.ID, &row.PostID, &row.ParentCommentID, &row.AuthorID, &row.Content, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("PsqlPool select comment by postID %w", err)
		}
		row.CreatedAt, row.UpdatedAt = t.parseTime()

		output = append(output, &row)
	}

	return output, nil
}

func (p PsqlPool) GetCommentByParentCommentID(ctx context.Context, parentCommentID string, first int32) ([]*model.Comment, error) {
	t := times{}

	var output []*model.Comment

	query := "SELECT * FROM comments WHERE parent_comment_id = $1 ORDER BY created_at DESC LIMIT 10 OFFSET $2"

	rows, err := p.Pool.Query(ctx, query, parentCommentID, first)
	defer rows.Close()

	if err != nil {
		return nil, fmt.Errorf("PsqlPool select comment by parentID %w", err)
	}

	var row model.Comment

	for rows.Next() {

		err = rows.Scan(&row.ID, &row.PostID, &row.ParentCommentID, &row.AuthorID, &row.Content, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("PsqlPool select comment by parentID %w", err)
		}
		row.CreatedAt, row.UpdatedAt = t.parseTime()

		output = append(output, &row)
	}

	return output, nil
}
