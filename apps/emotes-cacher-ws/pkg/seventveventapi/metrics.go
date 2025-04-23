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
)

type ClientMetrics struct {
	TotalShards int
	AliveShards int
	DeadShards  int
}
