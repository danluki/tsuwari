package bots

import (
	"sync"
	"time"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/chat_client"
	"github.com/satont/twir/apps/bots/pkg/tlds"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/giveaways"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"

	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/parser"

	model "github.com/satont/twir/libs/gomodels"

	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	DB             *gorm.DB
	Logger         logger.Logger
	Cfg            cfg.Config
	ParserGrpc     parser.ParserClient
	TokensGrpc     tokens.TokensClient
	EventsGrpc     events.EventsClient
	WebsocketsGrpc websockets.WebsocketClient
	GiveawaysGrpc  giveaways.GiveawaysClient

	Tlds  *tlds.TLDS
	Redis *redis.Client
}

type Service struct {
	Instances map[string]*chat_client.ChatClient

	db            *gorm.DB
	logger        logger.Logger
	cfg           cfg.Config
	parserGrpc    parser.ParserClient
	tokensGrpc    tokens.TokensClient
	eventsGrpc    events.EventsClient
	giveawaysGrpc giveaways.GiveawaysClient
}

func NewBotsService(opts Opts) *Service {
	service := &Service{
		Instances:     make(map[string]*chat_client.ChatClient),
		db:            opts.DB,
		logger:        opts.Logger,
		cfg:           opts.Cfg,
		parserGrpc:    opts.ParserGrpc,
		tokensGrpc:    opts.TokensGrpc,
		eventsGrpc:    opts.EventsGrpc,
		giveawaysGrpc: opts.GiveawaysGrpc,
	}
	mu := sync.Mutex{}

	var bots []model.Bots
	err := opts.DB.
		Preload("Token").
		Preload("Channels").
		Find(&bots).
		Error
	if err != nil {
		panic(err)
	}

	joinRateLimiter, _ := ratelimiting.NewSlidingWindow(15, 25*time.Second)

	for _, bot := range bots {
		bot := bot
		instance := newBot(
			ClientOpts{
				DB:              opts.DB,
				Cfg:             opts.Cfg,
				Logger:          opts.Logger,
				Model:           &bot,
				ParserGrpc:      opts.ParserGrpc,
				TokensGrpc:      opts.TokensGrpc,
				EventsGrpc:      opts.EventsGrpc,
				WebsocketsGrpc:  opts.WebsocketsGrpc,
				GiveawaysGrpc:   opts.GiveawaysGrpc,
				Redis:           opts.Redis,
				JoinRateLimiter: joinRateLimiter,
				Tlds:            opts.Tlds,
			},
		)

		mu.Lock()
		service.Instances[bot.ID] = instance
		mu.Unlock()
	}

	return service
}
