package graph

import (
	"context"
	"ozon/internal/transport/graph/model"
	"ozon/pkg/logger"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Service interface {
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

type Subscription interface {
	Subscribe(ctx context.Context, postId string) chan *model.Comment
	Unsubscribe(ctx context.Context, postId string, ch chan *model.Comment)
	Publish(ctx context.Context, comment *model.Comment)
	Check(postId string) bool
}

type Resolver struct {
	service      Service
	logs         logger.Logger
	subscription Subscription
}

func NewResolver(srv Service, logs logger.Logger, subscription Subscription) *Resolver {
	return &Resolver{
		service:      srv,
		logs:         logs,
		subscription: subscription,
	}
}
