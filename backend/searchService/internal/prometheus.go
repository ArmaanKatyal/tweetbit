package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetrics struct {
	prefix                             string
	httpTransactionTotal               *prometheus.CounterVec
	httpResponseTimeHistogram          *prometheus.HistogramVec
	kafkaTransactionTotal              *prometheus.CounterVec
	elasticSearchTransactionTotal      *prometheus.CounterVec
	elasticSearchResponseTimeHistogram *prometheus.HistogramVec
	CreateTweetResponseTimeHistogram   *prometheus.HistogramVec
	buckets                            []float64
}

const (
	Ok                  = "200"
	BadRequest          = "400"
	InternalServerError = "500"
	GET                 = "GET"
	POST                = "POST"
	Success             = "SUCCESS"
	Error               = "ERROR"
)

func InitPromMetrics(prefix string, buckets []float64) *PromMetrics {
	return &PromMetrics{
		prefix: prefix,
		httpTransactionTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_requests_total",
			Help: "total HTTP requests processed",
		}, []string{"code", "method"}),
		httpResponseTimeHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    prefix + "_response_time_seconds",
			Help:    "Histogram of response time for handler",
			Buckets: buckets,
		}, []string{"code", "method"}),
		kafkaTransactionTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_kafka_transactions_total",
			Help: "total Kafka transactions processed",
		}, []string{"topic"}),
		elasticSearchTransactionTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_elasticsearch_transactions_total",
			Help: "total Elasticsearch transactions processed",
		}, []string{"code", "index"}),
		elasticSearchResponseTimeHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    prefix + "_elasticsearch_response_time_seconds",
			Help:    "Histogram of response time for Elasticsearch",
			Buckets: buckets,
		}, []string{"code", "index"}),
		CreateTweetResponseTimeHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    prefix + "_create_tweet_response_time_seconds",
			Help:    "Histogram of response time for CreateTweet",
			Buckets: buckets,
		}, []string{"code"}),
		buckets: buckets,
	}
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (pm *PromMetrics) ObserveResponseTime(code string, method string, time float64) {
	pm.httpResponseTimeHistogram.WithLabelValues(code, method).Observe(time)
}

func (pm *PromMetrics) IncHttpTransaction(code string, method string) {
	pm.httpTransactionTotal.WithLabelValues(code, method).Inc()
}

func (pm *PromMetrics) IncKafkaTransaction(topic string) {
	pm.kafkaTransactionTotal.WithLabelValues(topic).Inc()
}

func (pm *PromMetrics) IncESTransaction(code string, index string) {
	pm.elasticSearchTransactionTotal.WithLabelValues(code, index).Inc()
}

func (pm *PromMetrics) ObserveESResponseTime(code string, index string, time float64) {
	pm.elasticSearchResponseTimeHistogram.WithLabelValues(code, index).Observe(time)
}
