package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromMetrics struct {
	prefix                           string
	httpTransactionTotal             *prometheus.CounterVec
	httpResponseTimeHistogram        *prometheus.HistogramVec
	CreateTweetResponseTimeHistogram *prometheus.HistogramVec
	buckets                          []float64
}

type MetricsInput struct {
	Code   string
	Method string
	Route  string
}

const (
	Ok                  = "200"
	BadRequest          = "400"
	InternalServerError = "500"
	GET                 = "GET"
	POST                = "POST"
	Success             = "SUCCESS"
	Error               = "ERROR"
	VerifyToken         = "verifyToken"
	VerifyApiKey        = "verifyApiKey"
)

func InitPromMetrics(prefix string, buckets []float64) *PromMetrics {
	return &PromMetrics{
		prefix: prefix,
		httpTransactionTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: prefix + "_requests_total",
			Help: "total HTTP requests processed",
		}, []string{"code", "method", "route"}),
		httpResponseTimeHistogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    prefix + "_response_time_seconds",
			Help:    "Histogram of response time for handler",
			Buckets: buckets,
		}, []string{"code", "method", "route"}),
		buckets: buckets,
	}
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (pm *PromMetrics) ObserveResponseTime(code string, method string, route string, time float64) {
	pm.httpResponseTimeHistogram.WithLabelValues(code, method, route).Observe(time)
}

func (pm *PromMetrics) IncHttpTransaction(code string, method string, route string) {
	pm.httpTransactionTotal.WithLabelValues(code, method, route).Inc()
}
