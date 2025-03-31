package service

import (
	"context"
	"fmt"
	"go.uber.org/mock/gomock"
	serviceMock "ozon/internal/service/mocks"
	"ozon/internal/transport/graph/model"
	"reflect"
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
		name                 string
		input                model.PostCommentInput
		postByIdMockBehavior getPostByIdBehavior
		want                 *model.Comment
		wantErr              bool
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
			wantErr: false,
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
			want:    nil,
			wantErr: true,
		},
		{
			name: "successful comment creation when post is not found",
			input: model.PostCommentInput{
				PostID:  "-1",
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
				PostID:  "-1",
				Content: "Test comment",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc, ctx := gomock.WithContext(context.Background(), t)
			repo := serviceMock.NewMockRepository(mc)

			repo.EXPECT().
				GetPostByID(ctx, tt.input.PostID).
				Return(tt.postByIdMockBehavior.output.post, tt.postByIdMockBehavior.output.err)

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

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.PostComment() = %v, want %v", got, tt.want)
			}
		})
	}
}
