package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// Logger es la instancia global del logger de zap
var Logger *zap.Logger

// InitLogger inicializa el logger según el modo producción o desarrollo
func InitLogger(production bool) error {
	var config zap.Config
	if production {
		// Configuración optimizada para producción
		config = zap.NewProductionConfig()
	} else {
		// Configuración para desarrollo con colores y más detalle
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Cambia la clave del tiempo a "timestamp" y usa formato ISO8601
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	// Construye el logger con la configuración especificada
	Logger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// GinLogger es un middleware para Gin que registra cada petición HTTP con información relevante
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Guarda el tiempo de inicio
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Procesa la petición
		c.Next()

		latency := time.Since(start) // Calcula la duración de la petición
		// Registra la información usando zap.Logger.Info con campos estructurados
		Logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		)
	}
}
