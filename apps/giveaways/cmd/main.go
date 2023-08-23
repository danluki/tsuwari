package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/ktr0731/dept/app"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/satont/twir/apps/giveaways/grpc_impl"
	"github.com/satont/twir/apps/giveaways/internal/types/services"

	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/giveaways"
	"github.com/satont/twir/libs/grpc/servers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	_, appCtxCancel := context.WithCancel(context.Background())

	z, _ := zap.NewDevelopment()
	logger := z.Sugar()

	cfg, err := config.New()
	if err != nil || cfg == nil {
		fmt.Println(err)
		panic("Cannot load config of application")
	}

	zap.ReplaceGlobals(z)

	if cfg.AppEnv != "development" {
		http.Handle("/metrics", promhttp.Handler())
		go http.ListenAndServe("0.0.0.0:3000", nil)
	}

	if cfg.SentryDsn != "" {
		sentry.Init(
			sentry.ClientOptions{
				Dsn:              cfg.SentryDsn,
				Environment:      cfg.AppEnv,
				Debug:            true,
				TracesSampleRate: 1.0,
			},
		)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Fatalln(err)
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)
	defer d.Close()

	redisConnOpts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		panic(err)
	}
	redisClient := redis.NewClient(redisConnOpts)
	defer redisClient.Close()
	redisClient.Conn()

	s := &services.Services{
		Config: cfg,
		Logger: logger,
		Gorm:   db,
		Redis:  redisClient,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servers.GIVEAWAYS_SERVER_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	giveaways.RegisterGiveawaysServer(grpcServer, grpc_impl.NewServer(s))
	go grpcServer.Serve(lis)
	defer grpcServer.GracefulStop()

	logger.Info("Giveaways microservice started")
	app.Run()

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	<-exitSignal
	logger.Info("Exiting")
	appCtxCancel()
}
