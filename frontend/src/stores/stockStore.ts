// Importamos defineStore de Pinia para crear el store
import { defineStore } from 'pinia'
// Importamos ref de Vue para crear variables reactivas
import { ref } from 'vue'
// Importamos tipos para mejorar la tipificación
import type { StockRecommendation, StockRecommendationHistory } from '@/types'
// Importamos funciones API para obtener datos del backend
import {
  fetchStocks as fetchStocksApi,
  fetchStockDetail,
  fetchTopStocks as fetchTopStocksApi,
  fetchTickers,
  fetchStockRecommendationHistory,
} from '@/api/stockApi'

// Definimos el store llamado 'stocks' usando la composición API
export const useStockStore = defineStore('stocks', () => {
  // Estado reactivo para lista de acciones (stocks)
  const stocks = ref<StockRecommendation[]>([])
  // Estado reactivo para la acción seleccionada actualmente
  const currentStock = ref<StockRecommendation | null>(null)
  // Estado para las acciones top recomendadas
  const topStocks = ref<StockRecommendation[]>([])
  // Estado para la lista de tickers (símbolos de acciones)
  const tickers = ref<string[]>([])
  // Estado para saber si se está cargando algo
  const loading = ref(false)
  // Estado para guardar mensajes de error
  const error = ref<string | null>(null)
  // Estado para saber si hay más datos para cargar
  const hasMore = ref(true)
  // Estado para controlar la página actual (paginación)
  const currentPage = ref(1)
  // Estado para el historial de recomendaciones de una acción
  const recommendationHistory = ref<StockRecommendationHistory[]>([])

  // NUEVO: Guardamos los filtros actuales usados para cargar acciones
  const currentFilters = ref({
    ticker: '',
    rating: '',
    limit: 20,
  })

  // Función para obtener acciones desde la API con filtros y paginación
  const fetchStocksAction = async (filter: {
    ticker: string
    rating: string
    limit: number
    page?: number
    append?: boolean // Nuevo parámetro para añadir datos o reemplazar
  }) => {
    try {
      loading.value = true
      error.value = null

      // Definimos la página y si agregamos datos o reemplazamos
      const page = filter.page || 1
      const append = filter.append || false

      // Si no estamos agregando, reiniciamos página y guardamos filtros
      if (!append) {
        currentPage.value = 1
        currentFilters.value = {
          ticker: filter.ticker,
          rating: filter.rating,
          limit: filter.limit,
        }
      } else {
        // Si agregamos, actualizamos la página
        currentPage.value = page
      }

      // Llamamos a la API para obtener las acciones
      const response = await fetchStocksApi(filter.ticker, filter.rating, filter.limit, page)

      // Si agregamos y la página es mayor que 1, juntamos datos nuevos con existentes
      if (append && page > 1) {
        stocks.value = [...stocks.value, ...response.data]
      } else {
        // Si no, reemplazamos los datos actuales
        stocks.value = response.data
      }

      // Actualizamos si hay más datos para cargar
      hasMore.value = response.hasMore
    } catch (err) {
      error.value = 'Failed to fetch stocks' // Mensaje simple de error
      console.error(err) // Log del error en consola
    } finally {
      loading.value = false // Siempre quitamos el loading al final
    }
  }

  // Función para obtener el historial de recomendaciones de una acción
  const fetchRecommendationHistory = async (ticker: string) => {
    try {
      loading.value = true
      error.value = null
      const data = await fetchStockRecommendationHistory(ticker)
      // Guardamos solo el primer elemento del historial (puede ajustarse)
      recommendationHistory.value = [data.data[0]]
    } catch (err) {
      error.value = 'Failed to fetch recommendation history'
      console.error(err)
      recommendationHistory.value = []
    } finally {
      loading.value = false
    }
  }

  // Función para cargar más acciones usando los filtros y página actuales
  const loadMoreStocks = async () => {
    try {
      loading.value = true
      error.value = null

      // Calculamos la siguiente página
      const nextPage = currentPage.value + 1

      // Llamamos a la API con filtros guardados y siguiente página
      const response = await fetchStocksApi(
        currentFilters.value.ticker,
        currentFilters.value.rating,
        currentFilters.value.limit,
        nextPage,
      )

      // Añadimos los nuevos datos al array existente
      stocks.value = [...stocks.value, ...response.data]
      hasMore.value = response.hasMore
      currentPage.value = nextPage // Actualizamos la página
    } catch (err) {
      error.value = 'Failed to load more stocks'
      console.error(err)
    } finally {
      loading.value = false
    }
  }

  // Alternativa: función que unifica fetch inicial y cargar más, según parámetro append
  const fetchStocksUnified = async (filter: {
    ticker: string
    rating: string
    limit: number
    page?: number
    append?: boolean
  }) => {
    try {
      loading.value = true
      error.value = null

      const page = filter.page || 1
      const append = filter.append || false

      if (!append) {
        currentFilters.value = {
          ticker: filter.ticker,
          rating: filter.rating,
          limit: filter.limit,
        }
        currentPage.value = 1
      }

      const response = await fetchStocksApi(filter.ticker, filter.rating, filter.limit, page)

      if (append && page > 1) {
        stocks.value = [...stocks.value, ...response.data]
        currentPage.value = page
      } else {
        stocks.value = response.data
        currentPage.value = 1
      }

      hasMore.value = response.hasMore
    } catch (err) {
      error.value = 'Failed to fetch stocks'
      console.error(err)
    } finally {
      loading.value = false
    }
  }

  // Función para obtener detalles de una acción específica por ticker
  const fetchStockDetailAction = async (ticker: string) => {
    try {
      console.log('Fetching stock detail for ticker:', loading.value, ticker)
      loading.value = true
      error.value = null
      const response = await fetchStockDetail(ticker)
      console.log('Stock detail response:', response)
      currentStock.value = response // Guardamos el detalle recibido
    } catch (err) {
      error.value = 'Failed to fetch stock details'
      console.error(err)
    } finally {
      loading.value = false
    }
  }

  // Función para obtener las acciones top recomendadas
  const fetchTopStocks = async () => {
    try {
      const response = await fetchTopStocksApi()
      topStocks.value = response
    } catch (err) {
      console.error('Failed to fetch top stocks', err)
    }
  }

  // Función para obtener la lista de tickers disponibles
  const fetchTickersAction = async () => {
    try {
      const response = await fetchTickers()
      if (Array.isArray(response)) {
        tickers.value = response
      } else {
        tickers.value = []
        console.error('fetchTickers did not return an array:', response)
      }
    } catch (err) {
      console.error('Failed to fetch tickers', err)
    }
  }

  // Exportamos todo lo que queremos que sea accesible desde componentes
  return {
    stocks,
    currentStock,
    topStocks,
    tickers,
    loading,
    error,
    hasMore,
    currentPage,
    recommendationHistory,
    currentFilters, // Exportamos filtros actuales para debug o uso externo
    fetchStocks: fetchStocksAction, // Función principal para traer stocks
    fetchStocksUnified, // Alternativa unificada para fetch/load more
    loadMoreStocks, // Función para cargar más acciones
    fetchStockDetail: fetchStockDetailAction, // Detalles de una acción
    fetchTopStocks, // Top stocks recomendadas
    fetchTickers: fetchTickersAction, // Lista de tickers
    fetchRecommendationHistory, // Historial de recomendaciones
  }
})
