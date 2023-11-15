package giveaways

import (
	"context"
	"errors"

	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/giveaways"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Giveaways struct {
	*impl_deps.Deps
}

func (c *Giveaways) convertEntity(entity *model.ChannelGiveaway) *giveaways.Giveaway {
	return &giveaways.Giveaway{
		Id:                        entity.ID,
		Description:               entity.Description,
		ChannelId:                 entity.ChannelID,
		EndAt:                     entity.EndAt.String(),
		EligibleUserGroups:        entity.EligibleUserGroups,
		FollowersLuck:             int32(entity.FollowersLuck),
		Keyword:                   entity.Keyword,
		MessagesCountLuck:         int32(entity.MessagesCountLuck),
		RandomNumberFrom:          int32(entity.RandomNumberFrom),
		RandomNumberTo:            int32(entity.RandomNumberTo),
		RequiredMinFollowTime:     int32(entity.RequiredMinFollowTime),
		RequiredMinMessages:       int32(entity.RequireMinMessages),
		RequiredMinSubscriberTier: int32(entity.RequireMinSubscriberTier),
		RequiredMinSubscriberTime: int32(entity.RequiredMinSubscriberTime),
		Type:                      string(entity.Type),
		WinnerRandomNumber:        int32(entity.WinnerRandomNumber),
		WinnerCount:               int32(entity.WinnersCount),
		RequiredMinWatchTime:      int32(entity.RequireMinMessages),
		SubscribersLuck:           int32(entity.SubscribersLuck),
		SubscribersTier1Luck:      int32(entity.SubscribersTier1Luck),
		SubscribersTier2Luck:      int32(entity.SubscribersTier2Luck),
		SubscribersTier3Luck:      int32(entity.SubscribersTier3Luck),
		StartAt:                   entity.StartAt.String(),
		IsRunning:                 entity.IsRunning,
		IsFinished:                entity.IsFinished,
	}
}

func (c *Giveaways) GiveawaysGetCurrent(
	ctx context.Context,
	_ *emptypb.Empty,
) (*giveaways.Giveaway, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	var dbGiveaway model.ChannelGiveaway

	err := c.Db.WithContext(ctx).
		Where(`"channel_id" = ? AND "is_finished" != ?`, dashboardId, true).
		First(&dbGiveaway).
		Error
	if err != nil {
		return nil, err
	}

	return c.convertEntity(&dbGiveaway), nil
}

func (c *Giveaways) GiveawaysGetParticipants(
	ctx context.Context,
	req *giveaways.GetParticipantsRequest,
) (*giveaways.GetParticipantsResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var participants []*model.ChannelGiveawayParticipant
	err := c.Db.WithContext(ctx).
		Where(`"giveaway_id" = ? AND "display_name" LIKE ? AND "channel_id" = ?`, req.GiveawayId, "%"+req.GetQuery()+"%", dashboardId).
		Find(&participants).
		Error
	if err != nil {
		return nil, err
	}

	var count int64
	err = c.Db.WithContext(ctx).Where(`"giveaway_id" = ?`, req.GiveawayId).Count(&count).Error
	if err != nil {
		return nil, err
	}

	var convertedPaticipants []*giveaways.Winner
	for _, participant := range participants {
		convertedPaticipants = append(convertedPaticipants, &giveaways.Winner{
			UserId:      participant.UserID,
			DisplayName: participant.DisplayName,
		})
	}

	return &giveaways.GetParticipantsResponse{
		Winners:    convertedPaticipants,
		TotalCount: count,
	}, nil
}

func (c *Giveaways) GiveawaysCreateOrGet(
	ctx context.Context,
	req *giveaways.CreateOrGetRequest,
) (*giveaways.Giveaway, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var dbGiveaway model.ChannelGiveaway
	err := c.Db.WithContext(ctx).
		Where(`"channel_id" = ? AND "is_finished" == ? AND "is_running" == ?`, dashboardId, false, true).
		First(&dbGiveaway).
		Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	entity := &model.ChannelGiveaway{
		ChannelID:                 dashboardId,
		RequiredMinSubscriberTime: int(req.RequiredMinSubscriberTime),
		RequireMinMessages:        int(req.RequiredMinMessages),
		EligibleUserGroups:        req.EligibleUserGroups,
		Description:               req.Description,
		Keyword:                   req.Keyword,
		FollowersLuck:             int(req.FollowersLuck),
		RequiredMinFollowTime:     int(req.RequiredMinFollowTime),
		RandomNumberTo:            int(req.RandomNumberTo),
		RandomNumberFrom:          int(req.RandomNumberFrom),
		MessagesCountLuck:         int(req.MessagesCountLuck),
		SubscribersLuck:           int(req.SubscribersLuck),
		SubscribersTier1Luck:      int(req.SubscribersTier1Luck),
		SubscribersTier2Luck:      int(req.SubscribersTier2Luck),
		SubscribersTier3Luck:      int(req.SubscribersTier3Luck),
		RequiredMinWatchTime:      int(req.RequiredMinWatchTime),
		RequireMinSubscriberTier:  int(req.RequiredMinSubscriberTier),
		WinnerRandomNumber:        int(req.WinnerRandomNumber),
		WinnersCount:              int(req.WinnerCount),
		Type:                      model.GiveawayType(req.Type),
	}
	err = c.Db.WithContext(ctx).Create(entity).Error
	if err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}

func (c *Giveaways) GiveawaysUpdate(
	ctx context.Context,
	req *giveaways.UpdateRequest,
) (*giveaways.Giveaway, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	dbGiveaway := model.ChannelGiveaway{}
	err := c.Db.WithContext(ctx).
		Where(`"channel_id" = ? AND "is_finished" != ?`, dashboardId, true).
		Find(&dbGiveaway).
		Error
	if err != nil {
		return nil, err
	}

	err = c.Db.WithContext(ctx).Model(&model.ChannelGiveaway{}).Updates(model.ChannelGiveaway{
		IsRunning:            req.GetIsRunning(),
		IsFinished:           req.GetIsFinished(),
		Description:          req.GetDescription(),
		Keyword:              req.GetKeyword(),
		FollowersLuck:        int(req.GetFollowersLuck()),
		RandomNumberTo:       int(req.GetRandomNumberTo()),
		RandomNumberFrom:     int(req.GetRandomNumberFrom()),
		RequiredMinWatchTime: int(req.GetRequiredMinWatchTime()),
		// RequiredMinFollowTime: int(req.GetRequiredMinFollowTime()),
		RequireMinMessages: int(req.GetRequiredMinMessages()),
		SubscribersLuck:    int(req.GetSubscribersLuck()),
		WinnersCount:       int(req.GetWinnersCount()),
		Type:               model.GiveawayType(req.GetType()),
	}).Error
	if err != nil {
		return nil, err
	}

	return c.convertEntity(&dbGiveaway), nil
}

func (c *Giveaways) GiveawaysDelete(
	ctx context.Context,
	req *giveaways.DeleteRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	err := c.Db.WithContext(ctx).
		Where(`"channel_id = ?" AND "id = ?"`, dashboardId, req.Id).
		Delete(&model.ChannelGiveaway{}).
		Error
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Giveaways) GiveawaysGetById(
	ctx context.Context,
	req *giveaways.GetByIdRequest,
) (*giveaways.Giveaway, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	var dbGiveaway model.ChannelGiveaway
	err := c.Db.WithContext(ctx).
		Where(`"channel_id = ?" AND "id = ?"`, dashboardId, req.Id).
		Group(`"id`).
		First(&dbGiveaway).
		Error
	if err != nil {
		return nil, err
	}

	return c.convertEntity(&dbGiveaway), nil
}
