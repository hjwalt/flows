package runtime_sarama

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var messageProcessedCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Help: "Flows messages processed per partition",
	Name: "flows_messages_received",
}, []string{"topic", "partition"})
