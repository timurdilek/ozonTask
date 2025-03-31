package http

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo"
	"github.com/vektah/gqlparser/v2/ast"
	"ozon/internal/Subscription"
	"ozon/internal/transport/graph"
	"ozon/pkg/logger"
)

type Handler struct {
	service graph.Service
	log     logger.Logger
	ps      graph.Subscription
}

func NewHandler(e *echo.Echo, service graph.Service, log logger.Logger) {
	handler := &Handler{
		service: service,
		log:     log,
		ps:      Subscription.New(),
	}

	e.POST("/query", handler.graphqlHandler())
	e.GET("/query", handler.graphqlHandler())
	e.GET("/", handler.playgroundHandler())
}

func (h *Handler) graphqlHandler() echo.HandlerFunc {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(h.service, h.log, h.ps)}))

	srv.AddTransport(transport.Websocket{})

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return func(c echo.Context) error {
		srv.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}

func (h *Handler) playgroundHandler() echo.HandlerFunc {
	srv := playground.Handler("GraphQL", "/query")

	return func(c echo.Context) error {
		srv.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}
