package chat_client

import (
	"context"

	"github.com/satont/twir/libs/grpc/generated/giveaways"
	"go.uber.org/zap"
)

func (c *ChatClient) handleGiveaways(
	msg Message,
	userBadges []string,
) {
	if msg.DbStream.ID == "" {
		return
	}
	defer func() {
		_, err := c.services.GiveawaysGrpc.TryProcessParticipant(
			context.Background(),
			&giveaways.TryProcessParticipantRequest{
				UserId:      msg.User.ID,
				MessageText: msg.Message,
				ChannelId:   msg.Channel.ID,
				DisplayName: msg.User.DisplayName,
			},
		)
		if err != nil {
			zap.S().Error(err)
			return
		}
	}()
}
