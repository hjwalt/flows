package metric

import (
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var initSync sync.Mutex

var promMetrics *Prometheus

func PrometheusRetry() Retry {
	return GetPrometheus()
}

func PrometheusProduce() Produce {
	return GetPrometheus()
}

func PrometheusConsume() Consume {
	return GetPrometheus()
}

func GetPrometheus() *Prometheus {
	if promMetrics != nil {
		return promMetrics
	}

	initSync.Lock()
	defer initSync.Unlock()

	promMetrics = &Prometheus{
		retry: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Help: "Retry count",
			Name: "flows_retry_count",
		}, []string{"topic", "partition"}),
		messagesProduced: promauto.NewCounterVec(prometheus.CounterOpts{
			Help: "Flows total messages produced from source topic",
			Name: "flows_messages_produced",
		}, []string{"topic", "partition"}),
		messagesProcessed: promauto.NewCounterVec(prometheus.CounterOpts{
			Help: "Flows messages processed per partition",
			Name: "flows_messages_received",
		}, []string{"topic", "partition"}),
	}

	return promMetrics
}

type Prometheus struct {
	retry             *prometheus.GaugeVec
	messagesProduced  *prometheus.CounterVec
	messagesProcessed *prometheus.CounterVec
}

func (p *Prometheus) RetryCount(topic string, partition int32, trycount int64) {
	if p == nil {
		return
	}
	p.retry.WithLabelValues(topic, fmt.Sprintf("%d", partition)).Set(float64(trycount))
}

func (p *Prometheus) MessagesProducedIncrement(topic string, partition int32, count int64) {
	if p == nil {
		return
	}
	p.messagesProduced.WithLabelValues(topic, fmt.Sprintf("%d", partition)).Add(float64(count))
}

func (p *Prometheus) MessagesProcessedIncrement(topic string, partition int32, count int64) {
	if p == nil {
		return
	}
	p.messagesProcessed.WithLabelValues(topic, fmt.Sprintf("%d", partition)).Add(float64(count))
}
