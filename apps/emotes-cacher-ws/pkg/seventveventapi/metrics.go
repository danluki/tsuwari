package seventveventapi

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalShards = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "seventveventapi_shards_total",
		Help: "Total number of 7TV WS shards",
	})
	aliveShards = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "seventveventapi_shards_alive",
		Help: "Number of alive 7TV WS shards",
	})
	deadShards = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "seventveventapi_shards_dead",
		Help: "Number of dead 7TV WS shards",
	})

	totalSubscriptions = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "seventv_total_subscriptions",
		Help: "Total number of active subscriptions",
	})

	subscriptionsPerShard = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "seventv_subscriptions_per_shard",
		Help: "Number of subscriptions per shard",
	}, []string{"session_id"})

	subscribedChannels = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "seventv_subscribed_channels",
		Help: "Number of subscribed channels per shard",
	}, []string{"session_id"})
)

type ClientMetrics struct {
	TotalShards int
	AliveShards int
	DeadShards  int
}

type SubscriptionMetrics struct {
	SessionID              string `json:"session_id"`
	SubscriptionLimit      int32  `json:"subscription_limit"`
	CurrentSubscriptions   int32  `json:"current_subscriptions"`
	HeartbeatIntervalMs    int32  `json:"heartbeat_interval_ms"`
	CurrentHeartbeatCycles int32  `json:"current_heartbeat_cycles"`
	Active                 bool   `json:"active"`
	SubscribedCount        int    `json:"subscribed_count"`
}
