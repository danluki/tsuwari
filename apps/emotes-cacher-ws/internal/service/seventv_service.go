package service

import (
	"context"
	"time"

	"github.com/twirapp/twir/emotes-cacher-ws/pkg/emotes"
	"github.com/twirapp/twir/emotes-cacher-ws/pkg/seventveventapi"
	"github.com/twirapp/twir/libs/repositories/channels"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	LC                 fx.Lifecycle
	Client             *seventveventapi.Client
	ChannelsRepository channels.Repository
}

type SevenTvService struct {
	client             *seventveventapi.Client
	channelsRepository channels.Repository
}

func New(opts Opts) *SevenTvService {
	service := &SevenTvService{
		client:             opts.Client,
		channelsRepository: opts.ChannelsRepository,
	}

	opts.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				channels, err := opts.ChannelsRepository.GetMany(
					context.Background(),
					channels.GetManyInput{},
				)
				if err != nil {
					panic(err)
				}

				for _, channel := range channels {
					emoteSets, err := emotes.GetChannelSevenTvEmotesSets(channel.ID)
					if err != nil {
						panic(err)
					}

					for _, set := range emoteSets {
						err := opts.Client.SubscribeToEmoteSetPatch(context.Background(), set)
						if err != nil {
							panic(err)
						}
					}
				}
			}()

			go func() {
				for {
					opts.Client.CollectMetrics()
					time.Sleep(5 * time.Second)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {

			return nil
		},
	})

	return service
}
