package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetrics struct {
	Prefix                     string
	ClientConnected            prometheus.Gauge
	HttpTransactionTotal       *prometheus.CounterVec
	HttpResponseTimeHistogram  *prometheus.HistogramVec
	KafkaTransactionTotal      *prometheus.CounterVec
	KafkaResponseTimeHistogram *prometheus.HistogramVec
	Buckets                    []float64
}

func InitPromMetrics(prefix string, buckets []float64) *PromMetrics {
	return &PromMetrics{
		Prefix: prefix,
		ClientConnected: promauto.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_client_connected",
			Help: "Number of active client connections",
		}),
		HttpTransactionTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_requests_total",
			Help: "total HTTP requests processed",
		}, []string{"code", "method"}),
		HttpResponseTimeHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    prefix + "_response_time_seconds",
			Help:    "Histogram of response time for handler",
			Buckets: buckets,
		}, []string{"code", "method"}),
		KafkaTransactionTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_kafka_transactions_total",
			Help: "total Kafka transactions processed",
		}, []string{"topic"}),
		KafkaResponseTimeHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    prefix + "_kafka_response_time_seconds",
			Help:    "Histogram of response time for Kafka transactions",
			Buckets: buckets,
		}, []string{"topic"}),
	}
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
