package stateless

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var retryGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Help: "Retry count",
	Name: "flows_retry_count",
}, []string{"topic", "partition"})
