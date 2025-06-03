package main

import (
	"api-stock/internal/config"
	"api-stock/internal/repository"
	"api-stock/internal/repository/api"
	"api-stock/internal/repository/cockroachdb"
	"api-stock/internal/service"
	"context"
	_ "database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Configuración
	cfg := config.Load()

	// Conexión a la base de datos
	db, err := cockroachdb.Connect(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Inicializar repositorios
	stockRepo := repository.NewStockRepository(db)
	apiClient := api.NewRecommendationClient(cfg.APIToken, cfg.APIBaseURL)

	// Inicializar servicio
	apiService := service.NewExternalAPIService(apiClient, stockRepo)

	// Canal para manejar señales de terminación
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Ticker para sincronización periódica
	ticker := time.NewTicker(cfg.WorkerInterval)
	defer ticker.Stop()

	// Sincronización inicial
	if err := apiService.SyncRecommendations(context.Background()); err != nil {
		log.Printf("Initial sync failed: %v", err)
	}

	log.Println("Worker started successfully")

	// Bucle principal del worker
	for {
		select {
		case <-ticker.C:
			log.Println("Starting incremental sync...")
			if err := apiService.IncrementalSync(context.Background()); err != nil {
				log.Printf("Incremental sync failed: %v", err)
			} else {
				log.Println("Incremental sync completed successfully")
			}

		case <-done:
			log.Println("Worker is shutting down...")
			return
		}
	}
}
