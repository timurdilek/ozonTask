package main

import (
	"context"
	"ozon/internal/app"
	"ozon/pkg/logger"
)

func main() {

	ctx := context.Background()
	cfg, err := app.LoadConfig()
	logger.InitLogger()
	if err != nil {
		panic(err)
	}

	a := app.New(ctx, cfg)
	app.Run(ctx, a)

}
