package model

import (
	"time"
)

type GiveawayType string

const (
	GiveawayTypeByKeyword      GiveawayType = "BY_KEYWORD"
	GiveawayTypeByRandomNumber GiveawayType = "BY_RANDOM_NUMBER"
)

type ChannelGiveaway struct {
	ID                        string       `gorm:"column:id;type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Description               string       `gorm:"column:description;not null"                                json:"description"`
	Type                      GiveawayType `gorm:"column:type;not null;default:BY_KEYWORD"                    json:"type"`
	ChannelID                 string       `gorm:"column:channel_id;not null"                                 json:"channel_id"`
	CreatedAt                 time.Time    `gorm:"column:created_at;not null;default:now()"                   json:"created_at"`
	StartAt                   time.Time    `gorm:"column:start_at;not null"                                   json:"start_at"`
	EndAt                     time.Time    `gorm:"column:end_at;not null"                                     json:"end_at"`
	ClosedAt                  time.Time    `gorm:"column:closed_at;not null"                                  json:"closed_at"`
	IsRunning                 bool         `gorm:"column:is_running;not null;default:false"                   json:"is_running"`
	IsFinished                bool         `gorm:"column:is_finished;not null;default:false"                  json:"is_finished"`
	RequiredMinWatchTime      int          `gorm:"column:required_min_watch_time"                             json:"required_min_watch_time"`
	RequiredMinFollowTime     int          `gorm:"column:required_min_follow_time"                            json:"required_min_follow_time"`
	RequireMinMessages        int          `gorm:"column:require_min_messages"                                json:"require_min_messages"`
	RequireMinSubscriberTier  int          `gorm:"column:require_min_subscriber_tier"                         json:"require_min_subscriber_tier"`
	RequiredMinSubscriberTime int          `gorm:"column:required_min_subscriber_time"                        json:"required_min_subscriber_time"`
	EligibleUserGroups        string       `gorm:"column:eligible_user_groups;not null"                       json:"eligible_user_groups"`
	Keyword                   string       `gorm:"column:keyword"                                             json:"keyword"`
	RandomNumberFrom          int          `gorm:"column:random_number_from"                                  json:"random_number_from"`
	RandomNumberTo            int          `gorm:"column:random_number_to"                                    json:"random_number_to"`
	WinnerRandomNumber        int          `gorm:"column:winner_random_number"                                json:"winner_random_number"`
	WinnersCount              int          `gorm:"column:winners_count;not null"                              json:"winners_count"`
	FollowersLuck             int          `gorm:"column:followers_luck;not null;default:0"                   json:"followers_luck"`
	SubscribersLuck           int          `gorm:"column:subscribers_luck;not null;default:0"                 json:"subscribers_luck"`
	SubscribersTier1Luck      int          `gorm:"column:subscribers_tier1_luck;not null;default:0"           json:"subscribers_tier1_luck"`
	SubscribersTier2Luck      int          `gorm:"column:subscribers_tier2_luck;not null;default:0"           json:"subscribers_tier2_luck"`
	SubscribersTier3Luck      int          `gorm:"column:subscribers_tier3_luck;not null;default:0"           json:"subscribers_tier3_luck"`
	MessagesCountLuck         int          `gorm:"column:messages_count_luck;not null;default:0"              json:"messages_count_luck"`
}

func (c *ChannelGiveaway) TableName() string {
	return "channels_giveaways"
}
