package main

import (
	"github.com/kataras/iris/v12/middleware/grpc"
	"github.com/satont/twir/apps/giveaways/internal/gorm"
	"github.com/satont/twir/apps/giveaways/internal/redis"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	twirsentry "github.com/satont/twir/libs/sentry"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.NopLogger,
		fx.Provide(
			cfg.NewFx,
			gorm.New,
			redis.New,
			twirsentry.NewFx(twirsentry.NewFxOpts{Service: "giveaways"}),
			logger.NewFx(
				logger.Opts{
					Service: "giveaways",
				},
			),
			grpc.New,
		),
		fx.Invoke(
			redis.New,
			gorm.New,
			grpc.New,
		),
	).Run()
}
