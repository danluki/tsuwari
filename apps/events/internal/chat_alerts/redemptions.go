package chat_alerts

import (
	"context"
	"strings"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/grpc/events"
)

func (c *ChatAlerts) redemption(
	ctx context.Context,
	settings model.ChatAlertsSettings,
	req *events.RedemptionCreatedMessage,
) error {
	if !settings.Redemptions.Enabled {
		return nil
	}

	if len(settings.Redemptions.Messages) == 0 {
		return nil
	}

	sample := lo.Sample(settings.Redemptions.Messages)

	text := sample.Text
	text = strings.ReplaceAll(text, "{user}", req.UserName)
	text = strings.ReplaceAll(text, "{reward}", req.RewardName)

	if text == "" {
		return nil
	}

	return c.bus.Bots.SendMessage.Publish(
		bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelId,
			Message:        text,
			SkipRateLimits: true,
		},
	)
}
