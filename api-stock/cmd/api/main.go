package main

import (
	_ "api-stock/docs"
	"api-stock/internal/config"
	httpservice "api-stock/internal/delivery/http"
	"api-stock/internal/repository"
	"api-stock/internal/repository/api"
	"api-stock/internal/repository/cockroachdb"
	"api-stock/internal/service"
	"api-stock/pkg/errors"
	"api-stock/pkg/logger"
	"api-stock/pkg/metrics"
	"context"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @title Stock Recommendation API
// @version 1.0
// @description API for stock recommendations with machine learning features
// @contact.name API Support
// @contact.email support@stockapi.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := logger.InitLogger(false); err != nil {
		zap.L().Fatal("Error al inicializar logger", zap.Error(err))
	}
	defer func() {
		if err := logger.Logger.Sync(); err != nil {
			logger.Logger.Error("Error al sincronizar logger", zap.Error(err))
		}
	}()
	// 1. Cargar configuración
	logger.Logger.Info("Cargando configuración...")
	cfg := config.Load()

	// 2. Inicializar logger estructurado
	logger.Logger.Info("Inicializando logger...")
	if err := logger.InitLogger(false); err != nil {
		zap.L().Fatal("Error al inicializar logger", zap.Error(err))
	}
	defer func() {
		if err := logger.Logger.Sync(); err != nil {
			zap.L().Error("Error al hacer sync del logger", zap.Error(err))
		}
	}()

	// 3. Inicializar métricas
	logger.Logger.Info("Inicializando sistema de métricas...")
	metrics.Init()

	// 4. Conectar a la base de datos
	logger.Logger.Info("Conectando a la base de datos...")
	db, err := cockroachdb.Connect(cfg.DBURL)
	if err != nil {
		logger.Logger.Fatal("Error al conectar a la base de datos", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Logger.Error("Error al cerrar conexión a la base de datos", zap.Error(err))
		}
	}()

	// 5. Ejecutar migraciones
	logger.Logger.Info("Ejecutando migraciones...")
	if err := repository.RunMigrations(db); err != nil {
		logger.Logger.Fatal("Error al ejecutar migraciones", zap.Error(err))
	}

	// 6. Inicializar repositorios
	logger.Logger.Info("Inicializando repositorios...")
	stockRepo := repository.NewStockRepository(db)
	apiClient := api.NewRecommendationClient(cfg.APIToken, cfg.APIBaseURL)

	// 7. Inicializar servicios
	logger.Logger.Info("Inicializando servicios...")
	stockService := service.NewStockService(stockRepo)
	recommendationService := service.NewRecommendationService(stockRepo)
	apiService := service.NewExternalAPIService(apiClient, stockRepo)

	// 8. Sincronización inicial de datos
	logger.Logger.Info("Ejecutando sincronización inicial con la API externa...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiService.SyncRecommendations(ctx); err != nil {
		logger.Logger.Error("Error durante la sincronización inicial", zap.Error(err))
	} else {
		logger.Logger.Info("Sincronización inicial completada")
	}

	// 9. Configurar el router HTTP
	logger.Logger.Info("Configurando router HTTP...")
	router := gin.New()

	// 10. Middlewares
	logger.Logger.Info("Configurando middlewares...")
	router.Use(
		gin.Recovery(),
		logger.GinLogger(),
		errors.ErrorHandler,
		metrics.PrometheusMiddleware(),
	)

	// 11. Configurar rutas
	logger.Logger.Info("Configurando rutas HTTP...")
	httpservice.SetupRoutes(router, stockService, recommendationService)

	// 12. Rutas adicionales
	router.GET("/metrics", gin.WrapH(metrics.Handler()))
	router.StaticFile("/swagger.yaml", "./docs/swagger.yaml")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 13. Configurar el servidor HTTP
	logger.Logger.Info("Inicializando servidor HTTP...")
	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
		IdleTimeout:  60 * time.Second,
	}

	// 14. Manejo de señales para apagado elegante
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// 15. Iniciar servidor en una goroutine
	go func() {
		logger.Logger.Info("Servidor iniciado",
			zap.String("puerto", cfg.HTTPPort),
			zap.String("entorno", cfg.Environment),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("Fallo al iniciar servidor", zap.Error(err))
		}
	}()

	// 16. Esperar señal de terminación
	<-done
	logger.Logger.Info("Señal de apagado recibida")

	// 18. Apagado con contexto
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelShutdown()

	logger.Logger.Info("Cerrando servidor HTTP...")
	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Logger.Error("Error al cerrar servidor", zap.Error(err))
	} else {
		logger.Logger.Info("Servidor cerrado correctamente")
	}
}
