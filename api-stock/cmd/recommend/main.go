package main

import (
	"api-stock/internal/config"
	"api-stock/internal/repository"
	"api-stock/internal/repository/cockroachdb"
	"api-stock/internal/service"
	"context"
	"fmt"
	"log"
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

	// Inicializar repositorio y servicio
	stockRepo := repository.NewStockRepository(db)
	recommendationService := service.NewRecommendationService(stockRepo)

	// Obtener las mejores recomendaciones
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	limit := 5 // Número de recomendaciones a mostrar
	best, err := recommendationService.GetBestStocks(ctx, limit)
	if err != nil {
		log.Fatalf("Error getting recommendations: %v", err)
	}

	// Mostrar resultados
	fmt.Println("\nTop stock recommendations:")
	for i, rec := range best {
		fmt.Printf("%d. %s (%s)\n", i+1, rec.Ticker, rec.Company)
		fmt.Printf("   Recommendation: %s -> %s\n", rec.RatingFrom, rec.RatingTo)
		fmt.Printf("   Broker: %s, Action: %s\n", rec.Brokerage, rec.Action)
		fmt.Printf("   Target: %s - %s, Date: %s\n\n",
			rec.TargetFrom, rec.TargetTo, rec.Time.Format("2006-01-02"))
	}

	// Mostrar acciones similares para la primera recomendación
	if len(best) > 0 {
		similar, err := recommendationService.FindSimilarStocks(ctx, best[0].Ticker, 3)
		if err != nil {
			log.Printf("Error finding similar stocks: %v", err)
		} else {
			fmt.Printf("\nStocks similar to %s:\n", best[0].Ticker)
			for _, s := range similar {
				fmt.Printf("- %s (similarity: %.2f)\n", s.Ticker, s.Similarity)
			}
		}
	}
}
