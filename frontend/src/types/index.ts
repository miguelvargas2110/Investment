// Tipo que representa una recomendación de acción
export type StockRecommendation = {
  ticker: string // Símbolo de la acción (ejemplo: "AAPL")
  company: string // Nombre de la empresa
  rating_from: string // Calificación inicial (ejemplo: "Hold")
  rating_to: string // Calificación final (ejemplo: "Buy")
  target_from: number // Precio objetivo mínimo recomendado
  target_to: number // Precio objetivo máximo recomendado
  brokerage: string // Nombre de la casa de bolsa o analista que da la recomendación
  action: string // Tipo de acción recomendada (ejemplo: "buy", "sell")
  time: string // Fecha o momento en que se hizo la recomendación (string ISO)
}

// Interfaz para la respuesta de la API que devuelve recomendaciones
export interface APIResponse {
  items: StockRecommendation[] // Lista de recomendaciones recibidas
  next_page: string // String que indica la siguiente página para paginación (puede ser URL o token)
}

// Historial de recomendaciones para una acción específica
export interface StockRecommendationHistory {
  time: string // Momento en que se registró esta recomendación
  target_from: number // Precio objetivo mínimo en este momento
  target_to: number // Precio objetivo máximo en este momento
}
