package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"ozon/pkg/logger"
)

// ConnectionDB connects to the postgres DB using
func ConnectionDB(ctx context.Context, DBConn string) *pgxpool.Pool {

	log := logger.GetLogger()

	pool, err := pgxpool.New(ctx, DBConn)
	if err = pool.Ping(ctx); err != nil {
		log.Fatal("database connection error!", zap.Error(err))
	}

	log.Info("connection to postgresql successful")
	return pool
}
