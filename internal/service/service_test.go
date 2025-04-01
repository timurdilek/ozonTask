package service

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	serviceMock "ozon/internal/service/mocks"
	"ozon/internal/transport/graph/model"
	"testing"
)

func TestService_PostComment(t *testing.T) {
	type getPostByIdResp struct {
		post *model.Post
		err  error
	}
	type getPostByIdBehavior struct {
		inputId string
		output  getPostByIdResp
	}

	tests := []struct {
		name                  string
		input                 model.PostCommentInput
		postByIdMockBehavior  getPostByIdBehavior
		want                  *model.Comment
		wantErr               bool
		shouldCallGetPostByID bool
	}{
		{
			name: "successful comment creation when post is found and comments are allowed",
			input: model.PostCommentInput{
				PostID:  "1",
				Content: "Test comment",
			},
			postByIdMockBehavior: getPostByIdBehavior{
				output: getPostByIdResp{
					post: &model.Post{
						ID:                 "1",
						AreCommentsAllowed: true,
					},
					err: nil,
				},
			},
			want: &model.Comment{
				ID:      "100",
				PostID:  "1",
				Content: "Test comment",
			},
			wantErr:               false,
			shouldCallGetPostByID: true,
		},
		{
			name: "error when comments are prohibited",
			input: model.PostCommentInput{
				PostID:  "2",
				Content: "Test comment",
			},
			postByIdMockBehavior: getPostByIdBehavior{
				output: getPostByIdResp{
					post: &model.Post{
						ID:                 "2",
						AreCommentsAllowed: false,
					},
					err: nil,
				},
			},
			want:                  nil,
			wantErr:               true,
			shouldCallGetPostByID: true,
		},
		{
			name: "successful comment creation when post is not found",
			input: model.PostCommentInput{
				PostID:  "1",
				Content: "Test comment",
			},
			postByIdMockBehavior: getPostByIdBehavior{
				output: getPostByIdResp{
					post: nil,
					err:  fmt.Errorf("post not found"),
				},
			},
			want: &model.Comment{
				ID:      "100",
				PostID:  "1",
				Content: "Test comment",
			},
			wantErr:               false,
			shouldCallGetPostByID: true,
		},
		{
			name: "adding a comment that is too long",
			input: model.PostCommentInput{
				PostID:  "1",
				Content: string(make([]byte, 2001)),
			},
			postByIdMockBehavior: getPostByIdBehavior{
				output: getPostByIdResp{
					post: nil,
					err:  ErrIncorrectPostLen,
				},
			},
			wantErr: true,
		},
		{
			name: "adding a comment with an empty content",
			input: model.PostCommentInput{
				PostID:  "1",
				Content: "",
			},
			postByIdMockBehavior: getPostByIdBehavior{
				output: getPostByIdResp{
					post: nil,
					err:  ErrIncorrectContentLen,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc, ctx := gomock.WithContext(context.Background(), t)
			repo := serviceMock.NewMockRepository(mc)

			if tt.shouldCallGetPostByID {
				repo.EXPECT().
					GetPostByID(ctx, tt.input.PostID).
					Return(tt.postByIdMockBehavior.output.post, tt.postByIdMockBehavior.output.err)
			}

			if !tt.wantErr && tt.want != nil {
				repo.EXPECT().
					PostComment(ctx, tt.input).
					Return(tt.want, nil)
			}
			s := &Service{
				repo: repo,
			}
			got, err := s.PostComment(ctx, tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.PostComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.Equal(t, got, tt.want) {
				t.Errorf("Service.PostComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_CreatePost(t *testing.T) {

	tests := []struct {
		name    string
		input   model.CreatePostInput
		want    *model.Post
		wantErr bool
	}{
		{
			name: "adding a post that is too long",
			input: model.CreatePostInput{
				Content: string(make([]byte, 10001)),
			},
			wantErr: true,
		},
		{
			name: "adding a post with an empty content",
			input: model.CreatePostInput{
				Content: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc, ctx := gomock.WithContext(context.Background(), t)
			repo := serviceMock.NewMockRepository(mc)

			s := &Service{
				repo: repo,
			}

			got, err := s.CreatePost(ctx, tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreatePost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.Equal(t, got, tt.want) {
				t.Errorf("Service.CreatePost() = %v, want %v", got, tt.want)
			}
		})
	}
}
