package app

import (
	"context"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"ozon/internal/server"
	"ozon/internal/service"
	"ozon/internal/transport/http"
	"syscall"

	"ozon/internal/repository"
	"ozon/internal/transport/graph/model"
	"ozon/pkg/logger"
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

type App struct {
	service    *service.Service
	repository Repository
}

func New(ctx context.Context, cfg Config) *App {

	storage, err := LoadConfig()
	log := logger.GetLogger()
	if err != nil {
		log.Fatal("No database has chosen")
	}

	a := &App{}

	switch storage.Storage.DB {
	case "postgres":
		log.Info("postgresql storage")
		a.repository = repository.NewPsql(ctx, parseDBConn(cfg.Postgres))
	case "in_memory":
		a.repository = repository.NewInMemoryRepo()
		log.Info("in_memory storage")
	default:
		log.Fatal("No database has chosen")
	}
	a.service = service.New(a.repository)
	return a
}

func Run(ctx context.Context, a *App) {

	log := logger.Logger{Logger: logger.GetLogger()}

	e := echo.New()

	service := service.New(a.repository)

	http.NewHandler(e, service, log)

	srv := server.New(e.Server.Handler)

	go func() {
		if err := srv.Run(ctx); err != nil {
			log.Fatal("failed to run server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	indication := <-quit

	log.Info("Gracefully stopping the server", zap.String("caught signal", indication.String()))

	if err := srv.Stop(); err != nil {
		log.Error("failed to stop server", zap.String("err", err.Error()))
	}

}
