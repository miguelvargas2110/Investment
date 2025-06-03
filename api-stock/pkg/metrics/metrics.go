package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

// Definición del contador para el total de solicitudes HTTP, segmentado por método, ruta y estado HTTP
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// Definición del histograma para la duración de las solicitudes HTTP, segmentado por ruta
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.3, 0.5, 1, 3, 5, 10},
		},
		[]string{"path"},
	)
)

// Handler devuelve el manejador HTTP estándar para exponer las métricas Prometheus
func Handler() http.Handler {
	return promhttp.Handler()
}

// PrometheusMiddleware es un middleware para Gin que mide y registra métricas de las solicitudes HTTP
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()        // Marca el inicio del procesamiento
		path := c.Request.URL.Path // Obtiene la ruta solicitada

		c.Next() // Procesa la siguiente función en la cadena de middleware/controlador

		status := strconv.Itoa(c.Writer.Status()) // Obtiene el código HTTP de respuesta como string
		duration := time.Since(start).Seconds()   // Calcula la duración en segundos

		// Incrementa el contador de solicitudes HTTP con las etiquetas método, ruta y estado
		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			path,
			status,
		).Inc()

		// Observa la duración de la solicitud en el histograma, etiquetado por ruta
		httpRequestDuration.WithLabelValues(path).Observe(duration)
	}
}

// Init registra las métricas definidas en el registro global de Prometheus
func Init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}
