package http

import (
	"api-stock/internal/domain" // Importa las interfaces de servicio del dominio
	"time"

	"github.com/gin-contrib/cors" // Middleware para manejo de CORS
	"github.com/gin-gonic/gin"    // Framework web Gin
)

// SetupRoutes configura todas las rutas HTTP de la aplicación.
func SetupRoutes(router *gin.Engine, stockService domain.StockService, recommendationService domain.RecommendationService) {
	// Middleware CORS para permitir solicitudes desde otros orígenes
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                      // Permitir solicitudes desde cualquier origen
		AllowMethods:     []string{"GET", "OPTIONS"},         // Métodos permitidos
		AllowHeaders:     []string{"Origin", "Content-Type"}, // Cabeceras permitidas
		ExposeHeaders:    []string{"Content-Length"},         // Cabeceras expuestas al cliente
		AllowCredentials: true,                               // Permitir cookies y credenciales
		MaxAge:           12 * time.Hour,                     // Tiempo de caché de la política CORS
	}))

	// Crea un nuevo handler pasando los servicios necesarios (inyección de dependencias)
	handler := NewStockHandler(stockService, recommendationService)

	// Agrupa las rutas bajo el prefijo /http/v1 (versión de la API)
	apiGroup := router.Group("/http/v1")
	{
		// Ruta de prueba para verificar si el servicio está activo
		apiGroup.GET("/health", handler.HealthCheck)

		// Agrupa rutas relacionadas con recomendaciones bajo /recommendations
		recGroup := apiGroup.Group("/recommendations")
		{
			recGroup.GET("", handler.GetRecommendations)          // Retorna todas las recomendaciones
			recGroup.GET("/tickers", handler.GetAvailableTickers) // Retorna todos los tickers disponibles
			recGroup.GET("/best", handler.GetBestRecommendations) // Retorna las mejores recomendaciones
		}
	}
}
