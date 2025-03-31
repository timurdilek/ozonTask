package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.70

import (
	"context"
	"net/http"
	"ozon/internal/transport/graph/model"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	if input.Content == "" {
		r.logs.Debug("invalid input arguments")
		return nil, &gqlerror.Error{
			Message: "invalid argument",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	if len(input.Content) > 10000 {
		r.logs.Debug("input content is too long")
		return nil, &gqlerror.Error{
			Message: "content is too long",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug("Creating post", zap.Any("input", input))

	post, err := r.service.CreatePost(ctx, model.CreatePostInput{
		Content:            input.Content,
		AreCommentsAllowed: input.AreCommentsAllowed,
		AuthorID:           input.AuthorID,
	})
	if err != nil {
		r.logs.Error("failed to create post", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to create post",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return post, nil
}

// PostComment is the resolver for the postComment field.
func (r *mutationResolver) PostComment(ctx context.Context, input model.PostCommentInput) (*model.Comment, error) {
	if len(input.Content) > 2000 {
		r.logs.Info("comment is too long")
		return nil, &gqlerror.Error{
			Message: "comment is too long",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	if input.Content == "" {
		r.logs.Info("invalid input arguments")
		return nil, &gqlerror.Error{
			Message: "invalid argument",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug("Creating comment", zap.Any("input", input))

	comment, err := r.service.PostComment(ctx, model.PostCommentInput{
		PostID:          input.PostID,
		ParentCommentID: input.ParentCommentID,
		Content:         input.Content,
		AuthorID:        input.AuthorID,
	})

	if err != nil {
		r.logs.Error("failed to create comment", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to create comment",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.subscription.Publish(ctx, comment)

	return comment, nil
}

// PutPost is the resolver for the putPost field.
func (r *mutationResolver) PutPost(ctx context.Context, input model.PutPostInput) (*model.Post, error) {
	if input.ID == "" {
		r.logs.Debug("invalid input arguments: missing post ID")
		return nil, &gqlerror.Error{
			Message: "invalid argument: missing post ID",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	if input.Content != nil && len(*input.Content) > 2000 {
		r.logs.Debug("input content is too long")
		return nil, &gqlerror.Error{
			Message: "content is too long",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug("Updating post", zap.Any("input", input))

	post, err := r.service.PutPost(ctx, input)
	if err != nil {
		r.logs.Error("failed to update post", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to update post",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return post, nil
}

// PutComment is the resolver for the putComment field.
func (r *mutationResolver) PutComment(ctx context.Context, input model.PutCommentInput) (*model.Comment, error) {
	if input.ID == "" {
		r.logs.Debug("invalid input arguments: missing comment ID")
		return nil, &gqlerror.Error{
			Message: "invalid argument: missing comment ID",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	if len(input.Content) > 2000 {
		r.logs.Info("comment is too long")
		return nil, &gqlerror.Error{
			Message: "comment is too long",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	if input.Content == "" {
		r.logs.Info("invalid input arguments: empty content")
		return nil, &gqlerror.Error{
			Message: "invalid argument: empty content",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug("Updating comment", zap.Any("input", input))

	comment, err := r.service.PutComment(ctx, input)
	if err != nil {
		r.logs.Error("failed to update comment", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to update comment",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return comment, nil
}

// DeletePost is the resolver for the deletePost field.
func (r *mutationResolver) DeletePost(ctx context.Context, id string) (bool, error) {
	if id == "" {
		r.logs.Debug("invalid input arguments: missing post ID")
		return false, &gqlerror.Error{
			Message: "invalid argument: missing post ID",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug("Deleting post", zap.String("id", id))

	success, err := r.service.DeletePost(ctx, id)
	if err != nil {
		r.logs.Error("failed to delete post", zap.String("err", err.Error()))
		return false, &gqlerror.Error{
			Message: "failed to delete post",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return success, nil
}

// DeleteComment is the resolver for the deleteComment field.
func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (bool, error) {
	if id == "" {
		r.logs.Debug("invalid input arguments: missing comment ID")
		return false, &gqlerror.Error{
			Message: "invalid argument: missing comment ID",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug("Deleting comment", zap.String("id", id))

	success, err := r.service.DeleteComment(ctx, id)
	if err != nil {
		r.logs.Error("failed to delete comment", zap.String("err", err.Error()))
		return false, &gqlerror.Error{
			Message: "failed to delete comment",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return success, nil
}

// GetPost is the resolver for the getPost field.
func (r *queryResolver) GetPost(ctx context.Context, first int32) ([]*model.Post, error) {
	if first < 0 {
		r.logs.Debug("invalid input arguments: first must be positive")
		return nil, &gqlerror.Error{
			Message: "invalid argument: first must be positive",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug("Fetching posts", zap.Int32("first", first))

	posts, err := r.service.GetPost(ctx, first)
	if err != nil {
		r.logs.Error("failed to fetch posts", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to fetch posts",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return posts, nil
}

// GetPostByID is the resolver for the getPostById field.
func (r *queryResolver) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	if id == "" {
		r.logs.Debug("invalid input arguments: missing post ID")
		return nil, &gqlerror.Error{
			Message: "invalid argument: missing post ID",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug("Fetching post by ID", zap.String("id", id))

	post, err := r.service.GetPostByID(ctx, id)
	if err != nil {
		r.logs.Error("failed to fetch post by ID", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to fetch post by ID",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return post, nil
}

// GetCommentByPostID is the resolver for the getCommentByPostId field.
func (r *queryResolver) GetCommentByPostID(ctx context.Context, postID string, first int32) ([]*model.Comment, error) {
	if postID == "" {
		r.logs.Debug("invalid input arguments: missing post ID")
		return nil, &gqlerror.Error{
			Message: "invalid argument: missing post ID",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	if first < 0 {
		r.logs.Debug("invalid input arguments: first must be positive")
		return nil, &gqlerror.Error{
			Message: "invalid argument: first must be positive",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug("Fetching comments by post ID", zap.String("postID", postID), zap.Int32("first", first))

	comments, err := r.service.GetCommentByPostID(ctx, postID, first)
	if err != nil {
		r.logs.Error("failed to fetch comments by post ID", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to fetch comments by post ID",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return comments, nil
}

// GetCommentByParentCommentID is the resolver for the getCommentByParentCommentId field.
func (r *queryResolver) GetCommentByParentCommentID(ctx context.Context, parentCommentID string, first int32) ([]*model.Comment, error) {
	if parentCommentID == "" {
		r.logs.Debug("invalid input arguments: missing parent comment ID")
		return nil, &gqlerror.Error{
			Message: "invalid argument: missing parent comment ID",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	if first < 0 {
		r.logs.Debug("invalid input arguments: first must be positive")
		return nil, &gqlerror.Error{
			Message: "invalid argument: first must be positive",
			Extensions: map[string]interface{}{
				"code": http.StatusBadRequest,
			},
		}
	}

	r.logs.Debug("Fetching comments by parent comment ID", zap.String("parentCommentID", parentCommentID), zap.Int32("first", first))

	comments, err := r.service.GetCommentByParentCommentID(ctx, parentCommentID, first)
	if err != nil {
		r.logs.Error("failed to fetch comments by parent comment ID", zap.String("err", err.Error()))
		return nil, &gqlerror.Error{
			Message: "failed to fetch comments by parent comment ID",
			Extensions: map[string]interface{}{
				"code": http.StatusInternalServerError,
			},
		}
	}

	return comments, nil
}

// SubscriptionForComment is the resolver for the subscriptionForComment field.
func (r *subscriptionResolver) SubscriptionForComment(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	if !r.subscription.Check(postID) {
		_, err := r.service.GetPostByID(ctx, postID)
		if err != nil {
			return nil, err
		}
	}

	r.logs.Debug("creating new subscription", zap.String("postId", postID))

	ch := r.subscription.Subscribe(ctx, postID)

	go func() {
		<-ctx.Done()
		r.logs.Debug("Unsubscribing from comments", zap.String("postId", postID))
		r.subscription.Unsubscribe(ctx, postID, ch)
	}()

	return ch, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
