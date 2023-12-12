package main

import (
	"github.com/satont/twir/apps/giveaways/internal/gorm"
	"github.com/satont/twir/apps/giveaways/internal/grpc"
	"github.com/satont/twir/apps/giveaways/internal/redis"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/clients"
	"github.com/satont/twir/libs/grpc/generated/tokens"
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
			func(config cfg.Config) tokens.TokensClient {
				return clients.NewTokens(config.AppEnv)
			},
			grpc.New,
		),
		fx.Invoke(
			redis.New,
			gorm.New,
			grpc.New,
		),
	).Run()
}
