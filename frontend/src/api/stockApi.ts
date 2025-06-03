// Importa el tipo `StockRecommendation` desde el archivo de tipos
import type { StockRecommendation } from '../types/index.js'

// Define la URL base para la API, obtenida desde las variables de entorno (útil para cambiar entre local/desarrollo/producción)
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/http/v1'

/**
 * fetchStocks - Obtiene una lista de recomendaciones de acciones desde la API con filtros opcionales.
 * @param ticker - símbolo bursátil para filtrar.
 * @param rating - calificación para filtrar (ej. "buy", "sell", etc.).
 * @param limit - número máximo de resultados por página.
 * @param page - número de página a solicitar.
 * @returns objeto con array de recomendaciones y si hay más resultados.
 */
export const fetchStocks = async (
  ticker: string,
  rating: string,
  limit: number,
  page: number,
): Promise<{ data: StockRecommendation[]; hasMore: boolean }> => {
  const params = new URLSearchParams()
  if (ticker) params.append('ticker', ticker) // Agrega el filtro por ticker si está definido
  if (rating) params.append('rating', rating) // Agrega el filtro por rating si está definido
  params.append('limit', limit.toString()) // Número de resultados por página
  params.append('page', page.toString()) // Número de página

  const response = await fetch(`${API_BASE_URL}/recommendations?${params.toString()}`)
  if (!response.ok) {
    throw new Error('Failed to fetch stocks') // Lanza un error si la respuesta HTTP no fue exitosa
  }
  return response.json() // Devuelve la respuesta como JSON
}

/**
 * fetchStockDetail - Obtiene el detalle de la recomendación más reciente para un símbolo bursátil.
 * @param ticker - símbolo bursátil.
 * @returns una única recomendación de acción.
 */
export const fetchStockDetail = async (ticker: string): Promise<StockRecommendation> => {
  const response = await fetch(`${API_BASE_URL}/recommendations?ticker=${ticker}&limit=1`)
  if (!response.ok) {
    throw new Error('Failed to fetch stock detail')
  }
  const data = await response.json()
  return data.data[0] // Devuelve solo el primer elemento (más reciente)
}

/**
 * fetchTopStocks - Obtiene las mejores recomendaciones de acciones.
 * @param limit - número máximo de resultados (por defecto 20).
 * @returns array de recomendaciones destacadas.
 */
export const fetchTopStocks = async (limit = 20): Promise<StockRecommendation[]> => {
  const response = await fetch(`${API_BASE_URL}/recommendations/best?limit=${limit}`)
  if (!response.ok) {
    throw new Error('Failed to fetch top stocks')
  }
  const data = await response.json()
  return data.best_recommendations // Devuelve solo el campo `best_recommendations` del JSON
}

/**
 * fetchTickers - Obtiene la lista de tickers disponibles.
 * @returns array de strings con los símbolos bursátiles.
 */
export const fetchTickers = async (): Promise<string[]> => {
  const response = await fetch(`${API_BASE_URL}/recommendations/tickers`)
  if (!response.ok) {
    throw new Error('Failed to fetch tickers')
  }
  return response.json()
}

/**
 * fetchStockRecommendationHistory - Obtiene el historial de recomendaciones (última) para un ticker.
 * @param ticker - símbolo bursátil.
 * @returns los datos JSON completos de la recomendación.
 */
export const fetchStockRecommendationHistory = async (ticker: string) => {
  const res = await fetch(`${API_BASE_URL}/recommendations?ticker=${ticker}&limit=1`)
  if (!res.ok) throw new Error('Error fetching stock recommendation history')
  return await res.json() // Devuelve los datos completos de la recomendación
}
