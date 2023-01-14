package stateless

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var messageProducedCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Help: "Flows total messages produced from source topic",
	Name: "flows_messages_produced",
}, []string{"topic", "partition"})
