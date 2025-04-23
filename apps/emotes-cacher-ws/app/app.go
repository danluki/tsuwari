package app

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/emotes-cacher-ws/internal/service"
	"github.com/twirapp/twir/emotes-cacher-ws/pkg/seventveventapi"
	"github.com/twirapp/twir/libs/baseapp"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsrepositorypgx "github.com/twirapp/twir/libs/repositories/channels/pgx"
	"github.com/twirapp/twir/libs/uptrace"
	"go.uber.org/fx"
)

const serviceName = "emotes-cacher-ws"

var App = fx.Module(
	serviceName,
	baseapp.CreateBaseApp(baseapp.Opts{AppName: serviceName}),
	fx.Provide(
		fx.Annotate(
			channelsrepositorypgx.NewFx,
			fx.As(new(channelsrepository.Repository)),
		),
		seventveventapi.NewFx,
	),
	fx.Invoke(
		func(config cfg.Config) {
			// if config.AppEnv != "development" {
			http.Handle("/metrics", promhttp.Handler())
			go http.ListenAndServe("0.0.0.0:3000", nil)
			// }
		},
		service.New,
		uptrace.NewFx(serviceName),
		func(l logger.Logger) {
			l.Info("Emotes Cacher WS started")
		},
	),
)
