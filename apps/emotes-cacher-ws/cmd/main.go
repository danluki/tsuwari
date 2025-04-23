package main

import (
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/emotes-cacher-ws/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		logger.FxDiOnlyErrors,
		app.App,
	).Run()
}
