package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/constants"
	"github.com/satont/twir/libs/grpc/generated/giveaways"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	LC     fx.Lifecycle
	Logger logger.Logger
	DB     *gorm.DB
	Redis  *redis.Client
}

func New(opts Opts) (giveaways.GiveawaysServer, error) {
	service := &Impl{
		DB:     opts.DB,
		Logger: opts.Logger,
		Redis:  opts.Redis,
	}

	grpcNetListener, err := net.Listen(
		"tcp",
		fmt.Sprintf("0.0.0.0:%d", constants.DISCORD_SERVER_PORT),
	)
	if err != nil {
		return nil, err
	}
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionAge: 1 * time.Minute,
			},
		),
	)

	giveaways.RegisterGiveawaysServer(grpcServer, service)

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go grpcServer.Serve(grpcNetListener)
				opts.Logger.Info("Grpc server is running")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				grpcServer.GracefulStop()
				return nil
			},
		},
	)

	return service, nil
}

type Impl struct {
	giveaways.UnimplementedGiveawaysServer

	Logger logger.Logger
	DB     *gorm.DB
	Redis  *redis.Client
}

func (c *Impl) TryProcessParticipant(
	ctx context.Context,
	req *giveaways.TryProcessParticipantRequest,
) (*emptypb.Empty, error) {
	giveaway := model.ChannelGiveaway{}

	err := c.DB.WithContext(ctx).
		Where(`"channel_id = ? AND "is_finished" = ? AND "is_running" = ?`, req.GetChannelId(), true, false).
		First(&giveaway).
		Error
	if err != nil {
		c.Logger.Error(
			"cannot get giveaway",
			slog.Any("err", err),
			slog.String("channelId", req.GetChannelId()),
			slog.String("userId", req.GetUserId()),
		)

		return nil, err
	}

	text := req.GetMessageText()
	if giveaway.Type == model.GiveawayTypeByKeyword {
		if !strings.Contains(text, giveaway.Keyword) {
			return &emptypb.Empty{}, nil
		}
	} else if giveaway.Type == model.GiveawayTypeByRandomNumber {
		numFromText, err := strconv.Atoi(text)
		if err != nil {
			return &emptypb.Empty{}, nil
		}

		if numFromText != giveaway.WinnerRandomNumber {
			return &emptypb.Empty{}, nil
		}
	}

	var dbUser model.Users
	err = c.DB.WithContext(ctx).Where(`"id" = ?`, req.GetUserId()).First(&dbUser).Error
	if err != nil {
		c.Logger.Error("Cannot get user", slog.Any("err", err))
		return nil, err
	}

	var roles []*model.ChannelRoleUser
	err = c.DB.WithContext(ctx).Where(`"userId" = ?`, dbUser.ID).Preload("Role").Find(&roles).Error
	if err != nil {
		c.Logger.Error("Cannot get user roles", slog.Any("err", err))
		return nil, err
	}

	var userStats model.UsersStats
	err = c.DB.Where(`"userId" = ? AND "channelId" = ?`, dbUser.ID, req.GetChannelId()).
		First(&userStats).
		Error
	if err != nil {
		c.Logger.Error("Cannot get user stats", slog.Any("err", err))
		return nil, err
	}

	// TODO: check if user has all roles, and all required stats
	if userStats.Messages < int32(giveaway.RequireMinMessages) {
		return &emptypb.Empty{}, nil
	}

	if userStats.Watched < int64(giveaway.RequiredMinWatchTime) {
		return &emptypb.Empty{}, nil
	}

	newParticipant := model.ChannelGiveawayParticipant{
		UserID:               dbUser.ID,
		DisplayName:          req.GetDisplayName(),
		IsSubscriber:         userStats.IsSubscriber,
		MessagesCount:        int(userStats.Messages),
		UserStatsWatchedTime: userStats.Watched,
	}

	err = c.DB.WithContext(ctx).Create(&newParticipant).Error
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Impl) ChooseWinner(
	ctx context.Context,
	req *giveaways.ChooseWinnerRequest,
) (*giveaways.ChooseWinnerResponse, error) {
	giveaway := model.ChannelGiveaway{}

	err := c.DB.WithContext(ctx).
		Where(`"id" = ? "isRunning" = ? AND "isFinished" = ?`, req.GetGiveawayId(), true, false).
		Find(&giveaway).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Cannot find giveaway with this id")
		}

		c.Logger.Error(
			"cannot get giveaway",
			slog.Any("err", err),
		)

		return nil, err
	}

	var participants []*model.ChannelGiveawayParticipant
	err = c.DB.WithContext(ctx).
		Where(`"giveaway_id" = ?`, req.GetGiveawayId()).
		Find(&participants).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Cannot find giveaway with this id")
		}

		return nil, err
	}

	processedParticipants := make([]*giveaways.SimpleWinner, 0, len(participants))
	for _, participant := range participants {
		countOfTimes := 0
		if participant.IsSubscriber {
			countOfTimes += giveaway.SubscribersLuck
		}
		if participant.IsFollower {
			countOfTimes += giveaway.FollowersLuck
		}
		//TODO: add more conditions

		for i := 0; i < countOfTimes; i++ {
			processedParticipants = append(processedParticipants, &giveaways.SimpleWinner{
				UserId:      participant.UserID,
				DisplayName: participant.DisplayName,
			})
		}
	}

	winners := make([]*giveaways.SimpleWinner, giveaway.WinnersCount)
	for i := 0; i < len(winners); i++ {
		randInd := rand.Intn(len(processedParticipants))
		winners[i] = processedParticipants[randInd]
		processedParticipants[randInd] = processedParticipants[len(processedParticipants)-1]
		processedParticipants = processedParticipants[:len(processedParticipants)-1]
	}

	c.Redis.Del(ctx, fmt.Sprintf("giveaway:%s", giveaway.ChannelID))

	return &giveaways.ChooseWinnerResponse{
		Winners: winners,
	}, nil
}
