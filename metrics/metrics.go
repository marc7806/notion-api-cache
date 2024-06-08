package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"strconv"
	"time"
)

var metricPath = "/metrics"
var requestCounterMetric = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "notion_api_cache_requests_total",
	Help: "Total number of Notion API http endpoint requests",
}, []string{"method", "path", "status"})
var requestDurationMetric = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "notion_api_cache_request_duration_seconds",
	Help: "Duration of Notion API http endpoint requests",
}, []string{"method", "path", "status"})
var numberOfCachedDocumentsMetric = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "notion_api_cache_documents_total",
	Help: "Total number of documents cached in the application",
})

func SetNumberOfTotalCachedDocuments(num int) {
	numberOfCachedDocumentsMetric.Set(float64(num))
}

func requestCounterMetricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == metricPath {
			c.Next()
			return
		}

		start := time.Now()

		c.Next()

		elapsedTime := float64(time.Since(start)) / float64(time.Second)

		labels := prometheus.Labels{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": strconv.Itoa(c.Writer.Status()),
		}

		requestDurationMetric.With(labels).Observe(elapsedTime)
		requestCounterMetric.With(labels).Inc()
	}
}

func prometheusHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func RegisterMetrics(router *gin.Engine) {
	log.Println("Initializing metrics")
	router.Use(requestCounterMetricMiddleware())
	router.GET(metricPath, prometheusHandler())
}
