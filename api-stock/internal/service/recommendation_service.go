package service

import (
	"api-stock/internal/domain"
	"context"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

// recommendationService implementa la lógica para obtener recomendaciones de acciones.
// Mantiene un repositorio para acceder a datos, pesos para el modelo, cachés y sincronización.
type recommendationService struct {
	repo         domain.StockRepository           // interfaz para acceder a la base de datos
	modelWeights domain.ModelWeights              // pesos usados para calcular scores de recomendaciones
	cache        map[string][]domain.SimilarStock // caché para resultados de acciones similares
	cacheMutex   sync.RWMutex                     // mutex para proteger acceso a cache

	bestStocksCache      []domain.StockRecommendation // caché para las mejores recomendaciones
	bestStocksCacheMutex sync.RWMutex                 // mutex para proteger la caché de mejores acciones
	bestStocksCacheTime  time.Time                    // timestamp de la última actualización del caché
	bestStocksCacheTTL   time.Duration                // tiempo de vida del caché para mejores acciones
}

// Constructor que inicializa el servicio con un repositorio y pesos predefinidos
func NewRecommendationService(repo domain.StockRepository) domain.RecommendationService {
	return &recommendationService{
		repo: repo,
		modelWeights: domain.ModelWeights{
			// pesos para diferentes tipos de acción
			ActionWeights: map[string]float64{
				"initiated":      2.5,
				"target raised":  3.2,
				"target lowered": -1.5,
				"reiterated":     1.8,
				"updated":        2.0,
				"maintained":     1.0,
				// se pueden agregar más casos aquí
			},
			// pesos para diferentes ratings
			RatingWeights: map[string]float64{
				"buy":          3.0,
				"comprar":      3.0,
				"outperform":   2.7,
				"superar":      2.8,
				"neutral":      1.0,
				"market":       0.5,
				"underperform": -1.5,
				"sell":         -2.5,
				"vender":       -2.5,
			},
			// pesos para brokers
			BrokerageWeights: map[string]float64{
				"goldman":        1.3,
				"morgan":         1.2,
				"jp":             1.2,
				"morgan stanley": 1.2,
				"jpmorgan":       1.2,
				"bmo":            1.1,
				"oppenheimer":    1.0,
				"mizuho":         0.9,
			},
			RecentnessWeight: 0.1, // peso para la recencia temporal de la recomendación
		},
		cache:              make(map[string][]domain.SimilarStock), // inicializa cache vacía
		bestStocksCacheTTL: 5 * time.Minute,                        // TTL de 5 minutos para caché de mejores acciones
	}
}

// GetBestStocks devuelve las mejores acciones recomendadas, respetando un límite y usando caché
func (s *recommendationService) GetBestStocks(ctx context.Context, limit int) ([]domain.StockRecommendation, error) {
	// Valida el límite: si es <= 0 o > 100, asigna 10 por defecto
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// Intenta usar caché con lectura protegida
	s.bestStocksCacheMutex.RLock()
	if time.Since(s.bestStocksCacheTime) < s.bestStocksCacheTTL && len(s.bestStocksCache) > 0 {
		cached := s.bestStocksCache
		s.bestStocksCacheMutex.RUnlock()
		// Si hay más en caché que el límite, corta el slice
		if len(cached) > limit {
			return cached[:limit], nil
		}
		return cached, nil
	}
	s.bestStocksCacheMutex.RUnlock()

	// Si no está en caché o expiró, consulta las recomendaciones recientes (últimos 30 días)
	recentRecs, err := s.repo.GetRecentRecommendations(ctx, 30*24*time.Hour)
	if err != nil {
		return nil, err
	}

	// Calcula scores para cada ticker basado en las recomendaciones
	scores := s.calculateScores(recentRecs)
	// Ordena los tickers por score descendente
	sorted := s.sortByScore(scores)

	// Obtiene las recomendaciones top con detalle
	best, err := s.getTopRecommendations(ctx, sorted, limit)
	if err != nil {
		return nil, err
	}

	// Actualiza caché con exclusión de escritura
	s.bestStocksCacheMutex.Lock()
	s.bestStocksCache = best
	s.bestStocksCacheTime = time.Now()
	s.bestStocksCacheMutex.Unlock()

	return best, nil
}

// calculateScores calcula un mapa ticker -> score promedio basado en recomendaciones
func (s *recommendationService) calculateScores(recommendations []domain.StockRecommendation) map[string]float64 {
	scores := make(map[string]float64) // acumuladores de score por ticker
	counts := make(map[string]int)     // cantidad de recomendaciones por ticker

	for _, rec := range recommendations {
		score := s.calculateScore(rec) // score individual para esta recomendación
		scores[rec.Ticker] += score    // acumula
		counts[rec.Ticker]++           // cuenta
	}

	// Divide acumulado entre número de recomendaciones para promedio
	for ticker := range scores {
		scores[ticker] /= float64(counts[ticker])
	}

	return scores
}

// calculateScore calcula el score para una recomendación individual basado en sus características
func (s *recommendationService) calculateScore(rec domain.StockRecommendation) float64 {
	// Extrae scores parciales para las diferentes características
	features := map[string]float64{
		"action":    s.getActionScore(rec.Action),
		"rating":    s.getRatingScore(rec.RatingTo),
		"brokerage": s.getBrokerageScore(rec.Brokerage),
		"recency":   s.getRecencyScore(rec.Time),
	}

	// Suma todos los valores
	score := 0.0
	for _, val := range features {
		score += val
	}
	// Promedia para obtener un único score
	return score / float64(len(features))
}

// normalize convierte un string a minúsculas y quita espacios en los extremos
func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// getActionScore obtiene el peso para una acción dada según el modelo
func (s *recommendationService) getActionScore(action string) float64 {
	actionNorm := normalize(action)
	for key, val := range s.modelWeights.ActionWeights {
		if strings.Contains(actionNorm, key) {
			return val
		}
	}
	return 0.0
}

// getRatingScore obtiene el peso para un rating dado según el modelo
func (s *recommendationService) getRatingScore(rating string) float64 {
	ratingNorm := normalize(rating)
	for key, val := range s.modelWeights.RatingWeights {
		if strings.Contains(ratingNorm, key) {
			return val
		}
	}
	return 0.0
}

// getBrokerageScore obtiene el peso para un broker dado según el modelo
func (s *recommendationService) getBrokerageScore(brokerage string) float64 {
	brokerNorm := normalize(brokerage)
	for key, val := range s.modelWeights.BrokerageWeights {
		if strings.Contains(brokerNorm, key) {
			return val
		}
	}
	return 0.0
}

// getRecencyScore calcula el peso según la recencia temporal de la recomendación usando decaimiento exponencial
func (s *recommendationService) getRecencyScore(recTime time.Time) float64 {
	hoursAgo := time.Since(recTime).Hours()
	if hoursAgo <= 0 {
		return s.modelWeights.RecentnessWeight
	}
	lambda := 0.05 // tasa de decaimiento para la recencia
	decay := math.Exp(-lambda * hoursAgo)
	return s.modelWeights.RecentnessWeight * decay
}

// sortByScore ordena un slice de estructuras ticker-score por score descendente
func (s *recommendationService) sortByScore(scores map[string]float64) []struct {
	Ticker string
	Score  float64
} {
	var sorted []struct {
		Ticker string
		Score  float64
	}
	for ticker, score := range scores {
		sorted = append(sorted, struct {
			Ticker string
			Score  float64
		}{ticker, score})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score > sorted[j].Score
	})
	return sorted
}

// getTopRecommendations obtiene detalles de recomendaciones para los tickers ordenados hasta un límite
func (s *recommendationService) getTopRecommendations(ctx context.Context, sorted []struct {
	Ticker string
	Score  float64
}, limit int) ([]domain.StockRecommendation, error) {
	var best []domain.StockRecommendation
	for i, item := range sorted {
		if i >= limit {
			break
		}
		// Obtiene recomendaciones para un ticker, pide solo la primera para mostrar
		recs, _, err := s.repo.GetRecommendations(ctx, item.Ticker, 1, limit)
		if err != nil {
			return nil, err
		}
		if len(recs) > 0 {
			best = append(best, recs[0])
		}
	}
	return best, nil
}

// FindSimilarStocks busca las k acciones más similares a una dada usando similitud coseno y cache
func (s *recommendationService) FindSimilarStocks(ctx context.Context, ticker string, k int) ([]domain.SimilarStock, error) {
	// Revisa cache con lectura protegida
	s.cacheMutex.RLock()
	if cached, exists := s.cache[ticker]; exists {
		s.cacheMutex.RUnlock()
		return cached, nil
	}
	s.cacheMutex.RUnlock()

	// Obtiene características vectoriales para el ticker objetivo
	targetFeatures, err := s.repo.GetStockFeatures(ctx, ticker)
	if err != nil {
		return nil, err
	}

	// Obtiene características para todas las acciones
	allStocks, err := s.repo.GetAllStockFeatures(ctx)
	if err != nil {
		return nil, err
	}

	type result struct {
		stock domain.SimilarStock
		err   error
	}
	ch := make(chan result, len(allStocks))

	// Para cada acción, calcula en paralelo la similitud con el objetivo
	for _, stock := range allStocks {
		stock := stock // captura variable para goroutine
		go func() {
			if stock.Ticker == ticker {
				ch <- result{} // ignora la misma acción
				return
			}
			sim := cosineSimilarity(targetFeatures, stock.Features)
			ch <- result{stock: domain.SimilarStock{Ticker: stock.Ticker, Similarity: sim}}
		}()
	}

	var similarities []domain.SimilarStock
	for i := 0; i < len(allStocks); i++ {
		res := <-ch
		if res.err == nil && res.stock.Ticker != "" {
			similarities = append(similarities, res.stock)
		}
	}
	close(ch)

	// Ordena descendente por similitud
	sort.Slice(similarities, func(i, j int) bool {
		return similarities[i].Similarity > similarities[j].Similarity
	})

	// Limita la cantidad a k
	if len(similarities) > k {
		similarities = similarities[:k]
	}

	// Guarda en cache con exclusión de escritura
	s.cacheMutex.Lock()
	s.cache[ticker] = similarities
	s.cacheMutex.Unlock()

	return similarities, nil
}

// cosineSimilarity calcula la similitud coseno entre dos vectores representados como mapas string->float64
func cosineSimilarity(a, b map[string]float64) float64 {
	dotProduct := 0.0
	magnitudeA := 0.0
	magnitudeB := 0.0

	for key := range a {
		if val, exists := b[key]; exists {
			dotProduct += a[key] * val
		}
		magnitudeA += a[key] * a[key]
	}
	for _, val := range b {
		magnitudeB += val * val
	}

	magnitudeA = math.Sqrt(magnitudeA)
	magnitudeB = math.Sqrt(magnitudeB)

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}
	return dotProduct / (magnitudeA * magnitudeB)
}
