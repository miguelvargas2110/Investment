package domain

import (
	"context"
	"time"
)

//////////////////////////////
// Interfaces del Repositorio
//////////////////////////////

// StockRepository define métodos para interactuar con la base de datos de recomendaciones de acciones.
type StockRepository interface {
	// Obtiene recomendaciones de acciones filtradas por ticker, con paginación.
	GetRecommendations(ctx context.Context, ticker string, page int, limit int) ([]StockRecommendation, int, error)

	// Obtiene una lista de todos los tickers disponibles en la base de datos.
	GetAvailableTickers(ctx context.Context) ([]string, error)

	// Inserta múltiples recomendaciones en la base de datos.
	InsertRecommendations(ctx context.Context, recommendations []StockRecommendation) error

	// Obtiene recomendaciones recientes según un umbral de tiempo (ej: últimas 24h).
	GetRecentRecommendations(ctx context.Context, since time.Duration) ([]StockRecommendation, error)

	// Obtiene la recomendación más reciente.
	GetLatestRecommendation(ctx context.Context) (*StockRecommendation, error)

	// Elimina todas las recomendaciones (útil para reiniciar datos).
	DeleteAllRecommendations(ctx context.Context) error

	// Obtiene los features vectoriales de una acción específica (para recomendaciones basadas en similitud).
	GetStockFeatures(ctx context.Context, ticker string) (map[string]float64, error)

	// Obtiene los features vectoriales de todas las acciones.
	GetAllStockFeatures(ctx context.Context) ([]struct {
		Ticker   string
		Features map[string]float64
	}, error)

	// Verifica la conexión a la base de datos (para health check).
	Ping(ctx context.Context) error
}

//////////////////////////////
// Interfaces para API externa
//////////////////////////////

// ExternalAPI representa un cliente que se comunica con una API externa.
type ExternalAPI interface {
	// Obtiene un conjunto de recomendaciones desde una API paginada.
	GetRecommendations(ctx context.Context, nextPage string) ([]StockRecommendation, string, error)

	// Obtiene todas las recomendaciones disponibles (sin paginar, si la API lo permite).
	GetAllRecommendations(ctx context.Context) ([]StockRecommendation, error)
}

//////////////////////////////
// Interfaces de Servicios
//////////////////////////////

// StockService expone operaciones disponibles para el frontend (UI/API REST).
type StockService interface {
	// Retorna recomendaciones para un ticker (con paginación).
	GetRecommendations(ctx context.Context, ticker string, page, limit int) ([]StockRecommendation, int, error)

	// Lista de tickers disponibles.
	GetAvailableTickers(ctx context.Context) ([]string, error)

	// Verifica el estado del sistema (ej. conectividad a DB).
	HealthCheck(ctx context.Context) error
}

// RecommendationService encapsula lógica de negocio para sugerencias de inversión.
type RecommendationService interface {
	// Obtiene las mejores acciones para invertir, basadas en criterios internos (ej. puntajes).
	GetBestStocks(ctx context.Context, limit int) ([]StockRecommendation, error)

	// Busca acciones similares a un ticker dado (basado en features vectoriales, KNN u otra heurística).
	FindSimilarStocks(ctx context.Context, ticker string, k int) ([]SimilarStock, error)
}

// ExternalAPIService encapsula la lógica de sincronización entre la API externa y la base de datos.
type ExternalAPIService interface {
	// Realiza una sincronización completa desde la API externa.
	SyncRecommendations(ctx context.Context) error

	// Realiza una sincronización incremental (nuevas páginas desde la última guardada).
	IncrementalSync(ctx context.Context) error
}
