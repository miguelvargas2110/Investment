package service

import (
	"api-stock/internal/domain"
	"context"
	"log"
	"time"
)

// externalAPIService implementa la interfaz domain.ExternalAPIService.
// Se encarga de sincronizar recomendaciones bursátiles entre la API externa y el repositorio local.
type externalAPIService struct {
	client     domain.ExternalAPI     // Cliente para consumir la API externa de recomendaciones
	repo       domain.StockRepository // Repositorio para almacenar y consultar recomendaciones en la base de datos
	maxRetries int                    // Número máximo de reintentos para llamadas fallidas (no usado en este código pero reservado)
}

// NewExternalAPIService crea una nueva instancia de externalAPIService con el cliente API y repositorio dados.
func NewExternalAPIService(client domain.ExternalAPI, repo domain.StockRepository) domain.ExternalAPIService {
	return &externalAPIService{
		client:     client,
		repo:       repo,
		maxRetries: 3, // Valor por defecto para reintentos (no utilizado aquí)
	}
}

// SyncRecommendations sincroniza todas las recomendaciones de la API externa y las guarda en el repositorio.
// Esto es una sincronización completa que reemplaza las recomendaciones existentes.
func (s *externalAPIService) SyncRecommendations(ctx context.Context) error {
	// Llama al cliente API para obtener todas las recomendaciones actuales
	recommendations, err := s.client.GetAllRecommendations(ctx)
	if err != nil {
		return err // Error al llamar la API externa
	}

	// Si no hay recomendaciones nuevas, solo loguea y termina
	if len(recommendations) == 0 {
		log.Println("No recommendations received from API")
		return nil
	}

	// Borra todas las recomendaciones actuales del repositorio para reemplazarlas por las nuevas
	if err := s.repo.DeleteAllRecommendations(ctx); err != nil {
		return err // Error al borrar datos antiguos
	}

	batchSize := 100 // Tamaño de lote para inserción en bulk para evitar querys demasiado grandes

	// Inserta las recomendaciones en lotes para evitar saturar la base o la API
	for i := 0; i < len(recommendations); i += batchSize {
		end := i + batchSize
		if end > len(recommendations) {
			end = len(recommendations) // Ajusta el índice final si queda menos de un batch completo
		}

		// Inserta el lote actual en el repositorio
		if err := s.repo.InsertRecommendations(ctx, recommendations[i:end]); err != nil {
			return err // Error al insertar el lote
		}
		// Pausa para no saturar la base de datos o API, evitando bloqueos o limitaciones
		time.Sleep(500 * time.Millisecond)
	}

	return nil // Todo fue exitoso
}

// IncrementalSync sincroniza solo las recomendaciones nuevas desde la última actualización guardada.
// Esto es útil para actualizar gradualmente sin borrar todo.
func (s *externalAPIService) IncrementalSync(ctx context.Context) error {
	// Obtiene la recomendación más reciente guardada para saber desde cuándo pedir novedades
	latestRec, err := s.repo.GetLatestRecommendation(ctx)
	if err != nil {
		return err
	}

	var newRecommendations []domain.StockRecommendation // Acumula nuevas recomendaciones para insertar
	nextPage := ""                                      // Para paginación de la API externa

	// Loop para obtener páginas de recomendaciones desde la API externa
	for {
		// Llama a la API para obtener recomendaciones y la siguiente página (si hay)
		recommendations, page, err := s.client.GetRecommendations(ctx, nextPage)
		if err != nil {
			return err // Error llamando API
		}

		// Recorre cada recomendación recibida
		for _, rec := range recommendations {
			// Si ya llegamos a una recomendación anterior o igual a la última guardada, detenemos la inserción
			if latestRec != nil && !rec.Time.After(latestRec.Time) {
				// Inserta todas las recomendaciones nuevas acumuladas hasta ahora
				if err := s.repo.InsertRecommendations(ctx, newRecommendations); err != nil {
					return err // Error al insertar nuevas recomendaciones
				}
				return nil // Sincronización incremental terminada exitosamente
			}
			// Si la recomendación es más reciente, la agregamos al batch para insertar
			newRecommendations = append(newRecommendations, rec)
		}

		// Si no hay más páginas, salimos del loop
		if page == "" {
			break
		}
		// Actualiza el token o parámetro para la siguiente página
		nextPage = page

		// Pausa para respetar límites o evitar saturar la API externa
		time.Sleep(1 * time.Second)
	}

	// Si quedaron recomendaciones nuevas sin insertar después de todas las páginas, las insertamos
	if len(newRecommendations) > 0 {
		return s.repo.InsertRecommendations(ctx, newRecommendations)
	}

	// Si no hay nuevas recomendaciones, simplemente termina sin error
	return nil
}
