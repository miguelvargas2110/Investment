package api

import (
	"api-stock/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// recommendationClient implementa la interfaz domain.ExternalAPI.
// Se encarga de realizar llamadas a una API externa para obtener recomendaciones bursátiles.
type recommendationClient struct {
	client  *http.Client
	apiKey  string
	baseURL string
}

// NewRecommendationClient crea una nueva instancia del cliente de recomendaciones externas.
// Recibe como parámetros la API key y la URL base de la API.
func NewRecommendationClient(apiKey, baseURL string) domain.ExternalAPI {
	return &recommendationClient{
		client:  &http.Client{Timeout: 30 * time.Second},
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// GetRecommendations obtiene un conjunto de recomendaciones desde la API externa.
// Puede aceptar un token de paginación `nextPage` para continuar desde la última página consultada.
func (rc *recommendationClient) GetRecommendations(ctx context.Context, nextPage string) ([]domain.StockRecommendation, string, error) {
	url := rc.baseURL
	if nextPage != "" {
		url = fmt.Sprintf("%s?next_page=%s", url, nextPage)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("error al crear la solicitud: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", rc.apiKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("error en la solicitud: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("código de estado no exitoso: %d", resp.StatusCode)
	}

	var apiResponse domain.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, "", fmt.Errorf("error al decodificar JSON: %v", err)
	}

	return apiResponse.Items, apiResponse.NextPage, nil
}

// GetAllRecommendations obtiene todas las recomendaciones de la API externa manejando la paginación.
// Continúa consultando hasta que no haya más páginas disponibles.
func (rc *recommendationClient) GetAllRecommendations(ctx context.Context) ([]domain.StockRecommendation, error) {
	var allItems []domain.StockRecommendation
	nextPage := ""

	for {
		items, newNextPage, err := rc.GetRecommendations(ctx, nextPage)
		if err != nil {
			return nil, fmt.Errorf("error obteniendo página: %v", err)
		}

		allItems = append(allItems, items...)

		if newNextPage == "" {
			break // No hay más páginas
		}

		nextPage = newNextPage
		time.Sleep(2 * time.Second) // Pausa para evitar rate limits
	}

	return allItems, nil
}
