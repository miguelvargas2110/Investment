package domain

import "time"

// StockRecommendation representa una recomendación de una acción con todos sus detalles.
// Esta estructura puede ser utilizada tanto para respuestas de API como para almacenamiento.
// @StockRecommendation
type StockRecommendation struct {
	// Símbolo del ticker de la acción
	Ticker string `json:"ticker" example:"AAPL"`
	// Precio objetivo inferior
	TargetFrom string `json:"target_from" example:"150.00"`
	// Precio objetivo superior
	TargetTo string `json:"target_to" example:"175.00"`
	// Nombre de la empresa
	Company string `json:"company" example:"Apple Inc."`
	// Acción tomada (por ejemplo: aumento, reducción)
	Action string `json:"action" example:"aumentado"`
	// Nombre de la firma de corretaje que emitió la recomendación
	Brokerage string `json:"brokerage" example:"Goldman Sachs"`
	// Calificación anterior
	RatingFrom string `json:"rating_from" example:"neutral"`
	// Nueva calificación
	RatingTo string `json:"rating_to" example:"comprar"`
	// Momento de la recomendación
	Time time.Time `json:"time" example:"2023-01-15T00:00:00Z"`
}

// APIResponse representa la estructura de respuesta de la API externa.
// Incluye una lista de recomendaciones y un token para la siguiente página.
// @APIResponse
type APIResponse struct {
	// Lista de recomendaciones
	Items []StockRecommendation `json:"items"`
	// Token para paginación
	NextPage string `json:"next_page"`
}

// SimilarStock representa una acción similar con una puntuación de similitud.
// Utilizada para recomendaciones de acciones similares.
// @SimilarStock
type SimilarStock struct {
	// Símbolo del ticker de la acción similar
	Ticker string `json:"ticker" example:"MSFT"`
	// Puntaje de similitud (de 0 a 1)
	Similarity float64 `json:"similarity" example:"0.85"`
}

// ModelWeights contiene los pesos usados en el modelo de recomendación.
// Estos pesos determinan la importancia relativa de cada atributo.
type ModelWeights struct {
	// Pesos asignados a diferentes tipos de acciones
	ActionWeights map[string]float64
	// Pesos asignados a diferentes calificaciones
	RatingWeights map[string]float64
	// Pesos asignados a diferentes firmas de corretaje
	BrokerageWeights map[string]float64
	// Peso asignado a la recencia de la recomendación
	RecentnessWeight float64
}
