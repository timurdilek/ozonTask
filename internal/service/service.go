package service

import (
	"context"
	"errors"
	"ozon/internal/transport/graph/model"
)

const (
	maxPostLen    = 10000
	maxCommentLen = 2000
)

type Repository interface {
	CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error)
	PostComment(ctx context.Context, input model.PostCommentInput) (*model.Comment, error)
	PutPost(ctx context.Context, input model.PutPostInput) (*model.Post, error)
	PutComment(ctx context.Context, input model.PutCommentInput) (*model.Comment, error)
	DeletePost(ctx context.Context, id string) (bool, error)
	DeleteComment(ctx context.Context, id string) (bool, error)
	GetPost(ctx context.Context, first int32) ([]*model.Post, error)
	GetPostByID(ctx context.Context, id string) (*model.Post, error)
	GetCommentByPostID(ctx context.Context, postID string, first int32) ([]*model.Comment, error)
	GetCommentByParentCommentID(ctx context.Context, parentCommentID string, first int32) ([]*model.Comment, error)
}

type Service struct {
	repo Repository
}

func New(repository Repository) *Service {
	return &Service{repo: repository}
}

func (s Service) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {

	if input.Content == "" {
		return nil, ErrIncorrectContentLen
	}

	if len(input.Content) > maxPostLen {
		return nil, ErrIncorrectPostLen
	}

	post, err := s.repo.CreatePost(ctx, input)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s Service) PostComment(ctx context.Context, input model.PostCommentInput) (*model.Comment, error) {

	if input.Content == "" {
		return nil, ErrIncorrectContentLen
	}

	if len(input.Content) > maxCommentLen {
		return nil, ErrIncorrectCommentLen
	}

	post, err := s.repo.GetPostByID(ctx, input.PostID)

	if err == nil && post != nil {
		if !post.AreCommentsAllowed {
			return nil, errors.New("not allowed")
		}
	}

	comment, err := s.repo.PostComment(ctx, input)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s Service) PutPost(ctx context.Context, input model.PutPostInput) (*model.Post, error) {
	post, err := s.repo.PutPost(ctx, input)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s Service) PutComment(ctx context.Context, input model.PutCommentInput) (*model.Comment, error) {
	comment, err := s.repo.PutComment(ctx, input)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s Service) DeletePost(ctx context.Context, id string) (bool, error) {
	ok, err := s.repo.DeletePost(ctx, id)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (s Service) DeleteComment(ctx context.Context, id string) (bool, error) {
	ok, err := s.repo.DeleteComment(ctx, id)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (s Service) GetPost(ctx context.Context, first int32) ([]*model.Post, error) {
	post, err := s.repo.GetPost(ctx, first)

	return post, err
}

func (s Service) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	post, err := s.repo.GetPostByID(ctx, id)

	return post, err
}

func (s Service) GetCommentByPostID(ctx context.Context, postID string, first int32) ([]*model.Comment, error) {
	comments, err := s.repo.GetCommentByPostID(ctx, postID, first)

	return comments, err
}

func (s Service) GetCommentByParentCommentID(ctx context.Context, parentCommentID string, first int32) ([]*model.Comment, error) {
	comments, err := s.repo.GetCommentByParentCommentID(ctx, parentCommentID, first)

	return comments, err
}
