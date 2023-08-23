package services

import (
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Services struct {
	Logger *zap.SugaredLogger
	Config *config.Config
	Redis  *redis.Client
	Gorm   *gorm.DB
}
