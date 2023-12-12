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

	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/constants"
	"github.com/satont/twir/libs/grpc/generated/giveaways"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/twitch"
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

	TokensGrpc tokens.TokensClient
	Cfg        cfg.Config
}

func New(opts Opts) (giveaways.GiveawaysServer, error) {
	service := &Impl{
		DB:         opts.DB,
		Logger:     opts.Logger,
		Redis:      opts.Redis,
		TokensGrpc: opts.TokensGrpc,
		Cfg:        opts.Cfg,
	}

	grpcNetListener, err := net.Listen(
		"tcp",
		fmt.Sprintf("0.0.0.0:%d", constants.GIVEAWAYS_SERVER_PORT),
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

	Logger     logger.Logger
	DB         *gorm.DB
	Redis      *redis.Client
	Cfg        cfg.Config
	TokensGrpc tokens.TokensClient
}

func (c *Impl) TryProcessParticipant(
	ctx context.Context,
	req *giveaways.TryProcessParticipantRequest,
) (*emptypb.Empty, error) {
	giveaway := model.ChannelGiveaway{}
	err := c.DB.WithContext(ctx).
		Where(`"channel_id" = ? AND "is_finished" = ? AND "is_running" = ?`, req.GetChannelId(), false, true).
		Find(&giveaway).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "Cannot find currently running giveaway")
		}
		c.Logger.Error(
			"cannot get giveaway",
			slog.Any("err", err),
			slog.String("channelId", req.GetChannelId()),
			slog.String("userId", req.GetUserId()),
		)

		return nil, err
	}

	if giveaway.ID == "" {
		return &emptypb.Empty{}, nil
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
	if userStats.Messages <= int32(giveaway.RequireMinMessages) {
		return &emptypb.Empty{}, nil
	}

	if userStats.Watched <= int64(giveaway.RequiredMinWatchTime) {
		return &emptypb.Empty{}, nil
	}

	//TODO: this doesnt work
	// isFollower, err := c.isFollower(ctx, req.GetChannelId(), dbUser.ID)
	// if err != nil {
	// 	c.Logger.Error("Cannot check if user is follower", slog.Any("err", err))
	// 	return nil, err
	// }

	newParticipant := model.ChannelGiveawayParticipant{
		GiveawayID:           giveaway.ID,
		UserID:               dbUser.ID,
		DisplayName:          req.GetDisplayName(),
		IsSubscriber:         userStats.IsSubscriber,
		IsModerator:          userStats.IsMod,
		IsVip:                userStats.IsVip,
		IsFollower:           true,
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
		Where(`"id" = ? AND "is_finished" = ?`, req.GetGiveawayId(), false).
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

	if len(participants) == 0 {
		return nil, status.Error(codes.Canceled, "No participants")
	}

	if len(participants) <= giveaway.WinnersCount {
		return nil, status.Error(
			codes.OutOfRange,
			"Participants count must be greater than winners count",
		)
	}

	for _, participant := range participants {
		err = c.DB.WithContext(ctx).
			Model(participant).
			Update("is_winner", false).
			Error
		if err != nil {
			return nil, err
		}
	}

	processedParticipants := make([]*giveaways.SimpleWinner, 0, len(participants))
	for _, participant := range participants {
		countOfTimes := 1
		if participant.IsSubscriber {
			countOfTimes += giveaway.SubscribersLuck
		}
		if participant.IsFollower {
			countOfTimes += giveaway.FollowersLuck
		}

		for i := 0; i < countOfTimes; i++ {
			processedParticipants = append(processedParticipants, &giveaways.SimpleWinner{
				UserId:      participant.UserID,
				DisplayName: participant.DisplayName,
			})
		}
	}

	winners := make([]*giveaways.SimpleWinner, giveaway.WinnersCount)
	err = c.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(winners); i++ {
			randInd := rand.Intn(len(processedParticipants))
			winners[i] = processedParticipants[randInd]
			processedParticipants[randInd] = processedParticipants[len(processedParticipants)-1]
			processedParticipants = processedParticipants[:len(processedParticipants)-1]

			err = c.DB.WithContext(ctx).
				Where(`"giveaway_id" = ? AND "user_id" = ?`, giveaway.ID, winners[i].UserId).
				Model(&model.ChannelGiveawayParticipant{}).
				Update("is_winner", true).Error
			if err != nil {
				return err
			}
		}

		err = c.DB.WithContext(ctx).
			Where(`"id" = ?`, giveaway.ID).
			Model(&model.ChannelGiveaway{}).
			Update("is_running", false).
			Error

		return err
	})
	if err != nil {
		return nil, err
	}

	return &giveaways.ChooseWinnerResponse{
		Winners: winners,
	}, nil
}

func (c *Impl) isFollower(ctx context.Context, channelId, userId string) (bool, error) {
	twitchClient, err := twitch.NewAppClientWithContext(
		ctx,
		c.Cfg,
		c.TokensGrpc,
	)
	if err != nil {
		return false, err
	}

	follow, err := twitchClient.GetChannelFollows(
		&helix.GetChannelFollowsParams{
			BroadcasterID: channelId,
			UserID:        userId,
			First:         0,
			After:         "",
		},
	)
	if err != nil {
		return false, err
	}

	if follow.ErrorMessage != "" {
		return false, errors.New("Cannot get data from twitch")
	}

	if len(follow.Data.Channels) == 0 {
		return false, nil
	}

	return true, nil
}
