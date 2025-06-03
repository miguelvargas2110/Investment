package service

import (
	"api-stock/internal/domain"
	"context"
)

// stockService implementa la interfaz domain.StockService
// y actúa como capa de servicio para manejar la lógica relacionada con acciones y recomendaciones.
type stockService struct {
	repo domain.StockRepository
}

// NewStockService es el constructor que recibe un repositorio y retorna una instancia de stockService.
func NewStockService(repo domain.StockRepository) domain.StockService {
	return &stockService{repo: repo}
}

// GetRecommendations obtiene recomendaciones para un ticker específico,
// paginando resultados según page y limit.
// Se validan los parámetros para evitar valores fuera de rango.
func (s *stockService) GetRecommendations(ctx context.Context, ticker string, page, limit int) ([]domain.StockRecommendation, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	// Delegamos la obtención de datos al repositorio
	return s.repo.GetRecommendations(ctx, ticker, page, limit)
}

// GetAvailableTickers retorna una lista con todos los tickers disponibles en el repositorio.
func (s *stockService) GetAvailableTickers(ctx context.Context) ([]string, error) {
	return s.repo.GetAvailableTickers(ctx)
}

// HealthCheck verifica el estado de la conexión con el repositorio.
func (s *stockService) HealthCheck(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
