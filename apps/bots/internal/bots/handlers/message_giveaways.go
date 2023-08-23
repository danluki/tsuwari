// package handlers

// import (
// 	"context"
// 	"strconv"

// 	model "github.com/satont/twir/libs/gomodels"
// 	"github.com/satont/twir/libs/grpc/generated/giveaways"
// )

// func (c *Handlers) handleGiveaways(
// 	msg *Message,
// 	userBadges []string,
// ) {
// 	giveaway := model.ChannelGiveaway{}
// 	err := c.db.Where(`"channelId" = ? AND "end_at" != ? AND "closed_at" != ?`, msg.Channel.ID, nil, nil).Find(&giveaway).Errir
// 	if err != nil {
// 		c.logger.Error(err)
// 		return
// 	}

// 	switch giveaway.Type {
// 	case model.ChannelGiveAwayTypeByKeyword:
// 		if msg.Message != giveaway.Keyword.String {
// 			return
// 		}
// 		c.giveawaysGrpc.HandleChatMessage(context.Background(), &giveaways.HandleChatMessageRequest{
// 			Text: ,
// 		})
// 		break
// 	case model.ChannelGiveAwayTypeByRandomNumber:
// 		num, err := strconv.Atoi(msg.Message)
// 		if err != nil {
// 			return
// 		}

// 		break
// 	}
// }
