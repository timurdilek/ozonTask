package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"ozon/internal/transport/graph/model"
	"ozon/pkg/logger"
	"sort"
	"sync"
	"time"
)

type InMemoryRepo struct {
	memory map[string]model.Post
	mu     *sync.Mutex
	logger *zap.Logger
}

func NewInMemoryRepo() *InMemoryRepo {
	log := logger.GetLogger()
	var inMemoryStorage = make(map[string]model.Post)
	return &InMemoryRepo{memory: inMemoryStorage, mu: &sync.Mutex{}, logger: log}
}

func (i InMemoryRepo) CreatePost(ctx context.Context, input model.CreatePostInput) (*model.Post, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	id := uuid.New().String()

	output := model.Post{
		ID:                 id,
		AuthorID:           input.AuthorID,
		Content:            input.Content,
		AreCommentsAllowed: input.AreCommentsAllowed,
		CreatedAt:          time.Now().Format(time.DateTime),
		UpdatedAt:          "",
	}

	i.memory[id] = output

	return &output, nil

}

func (i InMemoryRepo) PostComment(ctx context.Context, input model.PostCommentInput) (*model.Comment, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	id := uuid.New().String()

	output := model.Comment{
		ID:              id,
		PostID:          input.PostID,
		ParentCommentID: input.ParentCommentID,
		AuthorID:        input.AuthorID,
		Content:         input.Content,
		CreatedAt:       time.Now().Format(time.DateTime),
		UpdatedAt:       "",
	}

	post := i.memory[input.PostID]

	if input.ParentCommentID == nil {
		post.Comments = append(post.Comments, &output)
	} else {
		parentComment := findCommentByID(post.Comments, *input.ParentCommentID)
		if parentComment == nil {
			return nil, errors.New("parent comment with ID not found")
		}
		parentComment.Replies = append(parentComment.Replies, &output)
	}

	i.memory[post.ID] = post

	return &output, nil
}

func findCommentByID(comments []*model.Comment, id string) *model.Comment {

	for _, comment := range comments {
		if comment.ID == id {
			return comment
		}

		if found := findCommentByID(comment.Replies, id); found != nil {
			return found
		}
	}
	return nil
}

func (i InMemoryRepo) PutPost(ctx context.Context, input model.PutPostInput) (*model.Post, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	output, ok := i.memory[input.ID]
	if !ok {
		return nil, errors.New("post does not exist")
	}
	output.CreatedAt = time.Now().Format(time.DateTime)

	if input.Content != nil {
		output.Content = *input.Content
	}
	if input.AreCommentsAllowed != nil {
		output.AreCommentsAllowed = *input.AreCommentsAllowed
	}

	return &output, nil
}

func (i InMemoryRepo) PutComment(ctx context.Context, input model.PutCommentInput) (*model.Comment, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	post, ok := i.memory[input.ID]
	if !ok {
		return nil, errors.New("post does not exist")
	}

	var output = findCommentByID(post.Comments, input.ID)
	if output == nil {
		return output, errors.New("there is no comment with this id")
	}

	output.Content = input.Content
	output.UpdatedAt = time.Now().Format(time.DateTime)

	return output, nil
}

func (i InMemoryRepo) DeletePost(ctx context.Context, id string) (bool, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	_, ok := i.memory[id]
	if !ok {
		return false, errors.New("post does not exist")
	}

	delete(i.memory, id)
	return true, nil
}

func (i InMemoryRepo) DeleteComment(ctx context.Context, id string) (bool, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	var deleteComment func(comments []*model.Comment) bool
	deleteComment = func(comments []*model.Comment) bool {
		for i := 0; i < len(comments); i++ {
			if comments[i].ID == id {
				comments = append(comments[:i], comments[i+1:]...)
				return true
			}
			if len(comments[i].Replies) > 0 {
				if deleted := deleteComment(comments[i].Replies); deleted {
					return true
				}
			}
		}
		return false
	}

	for _, post := range i.memory {
		if deleted := deleteComment(post.Comments); deleted {
			return true, nil
		}
	}

	return false, nil
}

func (i InMemoryRepo) GetPost(ctx context.Context, first int32) ([]*model.Post, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if first < 0 {
		return nil, errors.New("invalid 'first' value")
	}
	var posts []*model.Post

	for _, post := range i.memory {
		posts = append(posts, &post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID < posts[j].ID
	})

	start := int(first)
	end := start + 10

	if start >= len(posts) {
		return nil, nil
	}

	if end > len(posts) {
		end = len(posts)
	}

	return posts[start:end], nil
}

func (i InMemoryRepo) GetPostByID(ctx context.Context, id string) (*model.Post, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	post, exists := i.memory[id]
	if !exists {
		return nil, errors.New("post with this ID not found")
	}

	return &post, nil
}

func (i InMemoryRepo) GetCommentByPostID(ctx context.Context, postID string, first int32) ([]*model.Comment, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	post, exists := i.memory[postID]
	if !exists {
		return nil, errors.New("post with this ID not found")
	}

	return post.Comments, nil
}

func (i InMemoryRepo) GetCommentByParentCommentID(ctx context.Context, parentCommentID string, first int32) ([]*model.Comment, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	var result []*model.Comment

	for _, post := range i.memory {
		result = append(result, findCommentsByParentID(post.Comments, parentCommentID)...)
	}

	if first > 0 && int(first) < len(result) {
		result = result[:first]
	}

	return result, nil
}

func findCommentsByParentID(comments []*model.Comment, parentCommentID string) []*model.Comment {
	var result []*model.Comment

	for _, comment := range comments {
		if comment.ParentCommentID != nil && *comment.ParentCommentID == parentCommentID {
			result = append(result, comment)
		}
		result = append(result, findCommentsByParentID(comment.Replies, parentCommentID)...)
	}

	return result
}
