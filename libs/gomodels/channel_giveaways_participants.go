package model

import (
	"time"
)

type ChannelGiveawayParticipant struct {
	ID                   string    `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	GiveawayID           string    `gorm:"column:giveaway_id;not null"                                json:"giveawayId"`
	IsWinner             bool      `gorm:"column:is_winner;not null;default:false"                    json:"isWinner"`
	UserID               string    `gorm:"column:user_id;not null"                                    json:"userId"`
	DisplayName          string    `gorm:"column:display_name;not null"                               json:"displayName"`
	IsSubscriber         bool      `gorm:"column:is_subscriber;not null;default:false"                json:"isSubscriber"`
	IsFollower           bool      `gorm:"column:is_follower;not null;default:false"                  json:"isFollower"`
	IsModerator          bool      `gorm:"column:is_moderator;not null;default:false"                 json:"isModerator"`
	IsVip                bool      `gorm:"column:is_vip;not null;default:false"                       json:"isVip"`
	SubscriberTier       int       `gorm:"column:subscriber_tier"                                     json:"subscriberTier"`
	UserFollowSince      time.Time `gorm:"column:user_follow_since"                                   json:"userFollowSince"`
	UserStatsWatchedTime int64     `gorm:"column:user_stats_watched_time;not null"                    json:"userStatsWatchedTime"`
	MessagesCount        int       `gorm:"column:messages_count;not null;default:0"                   json:"messagesCount"`
}

func (c *ChannelGiveawayParticipant) TableName() string {
	return "channels_giveaways_participants"
}
