<template>
  <div class="h-full flex flex-col">
    <!-- Sección de filtros para búsqueda y filtrado de stocks -->
    <div class="bg-white rounded-xl shadow-sm p-6 mb-4 border border-gray-100">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <!-- Input para filtrar por ticker -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Ticker</label>
          <input
            v-model="filter.ticker"
            type="text"
            placeholder="Search ticker..."
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <!-- Select para filtrar por rating -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Rating</label>
          <select
            v-model="filter.rating"
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">All Ratings</option>
            <!-- Se listan las ratings disponibles dinámicamente -->
            <option v-for="rating in availableRatings" :key="rating" :value="rating">
              {{ rating }}
            </option>
          </select>
        </div>

        <!-- Input para filtrar por brokerage -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Brokerage</label>
          <input
            v-model="filter.brokerage"
            type="text"
            placeholder="Filter by brokerage..."
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <!-- Input para limitar la cantidad de resultados mostrados -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Items per page</label>
          <input
            v-model="filter.limit"
            type="number"
            placeholder="Number of results..."
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <!-- Botón para resetear los filtros a valores iniciales -->
        <div class="md:col-span-4 flex items-end w-full">
          <button
            @click="resetFilters"
            class="w-full bg-blue-200 hover:bg-gray-300 text-gray-800 py-2 px-4 rounded-md transition-colors"
          >
            Reset Filters
          </button>
        </div>
      </div>
    </div>

    <!-- Estado de carga principal, se muestra un spinner mientras loading es true y no está cargando más -->
    <div
      v-if="loading && !loadingMore"
      class="flex-1 flex flex-col items-center justify-center py-12 space-y-4"
    >
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-600"></div>
      <p class="text-gray-600 font-medium">Loading recommendations...</p>
    </div>

    <!-- Estado de error, se muestra si hay error -->
    <div
      v-else-if="error"
      class="flex-1 bg-red-50 border-l-4 border-red-400 p-4 mb-6 flex items-center"
    >
      <div class="flex">
        <div class="flex-shrink-0">
          <!-- Icono de error -->
          <svg
            class="h-5 w-5 text-red-400"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
              clip-rule="evenodd"
            />
          </svg>
        </div>
        <div class="ml-3">
          <p class="text-sm text-red-700">{{ error }}</p>
          <!-- Mensaje de error -->
          <button @click="$emit('refresh')" class="mt-2 text-sm text-blue-600 hover:text-blue-800">
            Retry
            <!-- Botón para reintentar, emite evento 'refresh' -->
          </button>
        </div>
      </div>
    </div>

    <!-- Contenedor principal de la tabla de stocks -->
    <div v-else class="flex-1 flex flex-col min-h-0">
      <div
        class="bg-white shadow-sm rounded-xl overflow-hidden border border-gray-200 flex-1 flex flex-col"
      >
        <div class="overflow-auto flex-1">
          <!-- Tabla de recomendaciones -->
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50 sticky top-0">
              <tr>
                <!-- Encabezados de la tabla -->
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Ticker
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Company
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Brokerage
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Action
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Rating
                </th>
                <th
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Target Price
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <!-- Renderizado de filas con stocks filtrados -->
              <tr
                v-for="(stock, index) in filteredStocks"
                :key="index"
                class="hover:bg-gray-50 transition-colors"
              >
                <!-- Datos de cada columna -->
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="font-medium text-gray-900">{{ stock.ticker }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-gray-900">{{ stock.company }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-gray-900">{{ stock.brokerage }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    :class="getActionClass(stock.action)"
                    class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                  >
                    {{ formatAction(stock.action) }}
                    <!-- Texto formateado de la acción -->
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    :class="getRatingClass(stock.ratingTo)"
                    class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                  >
                    {{ stock.ratingTo }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <!-- Si el targetFrom y targetTo son distintos, muestra el anterior tachado -->
                    <span
                      v-if="stock.targetFrom !== stock.targetTo"
                      class="text-gray-500 line-through mr-2"
                    >
                      {{ formatCurrency(stock.targetFrom) }}
                    </span>
                    <!-- Muestra el targetTo con clase según cambio -->
                    <span :class="getTargetChangeClass(stock.targetFrom, stock.targetTo)">
                      {{ formatCurrency(stock.targetTo) }}
                    </span>
                    <!-- Componente para mostrar indicador gráfico de cambio -->
                    <PriceChangeIndicator :from="stock.targetFrom" :to="stock.targetTo" />
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Estado vacío cuando no hay stocks filtrados -->
        <div
          v-if="filteredStocks.length === 0 && !loadingMore"
          class="flex-1 flex items-center justify-center text-center py-12"
        >
          <div>
            <svg
              class="mx-auto h-12 w-12 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="1.5"
                d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <h3 class="mt-4 text-lg font-medium text-gray-900">No recommendations found</h3>
            <p class="mt-2 text-gray-500">Try adjusting your search filters</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Botón para cargar más stocks, visible si hay más por cargar y no está cargando ahora -->
    <div v-if="hasMore && !loading" class="flex justify-center mt-6 px-4">
      <button
        @click="handleLoadMore"
        class="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg shadow-md transition-all duration-200 flex items-center w-full md:w-auto"
        :disabled="loadingMore"
        :class="{
          'opacity-75 cursor-not-allowed': loadingMore,
          'hover:shadow-lg': !loadingMore,
        }"
      >
        <!-- Spinner mientras carga más -->
        <svg
          v-if="loadingMore"
          class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle
            class="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            stroke-width="4"
          ></circle>
          <path
            class="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          ></path>
        </svg>
        <!-- Texto cambia según estado de carga -->
        {{ loadingMore ? 'Loading more stocks...' : 'Load More Stocks' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
// Importa componente para mostrar indicador gráfico de cambio de precio
import PriceChangeIndicator from '@/components/stock/PriceChangeIndicator.vue'

// Define la interfaz para los objetos stock recibidos
interface Stock {
  ticker: string
  company: string
  brokerage: string
  action: string
  ratingTo: string
  targetFrom: string
  targetTo: string
}

// Props para recibir datos y estado desde componente padre
const props = defineProps<{
  stocks: Stock[]
  loading: boolean
  loadingMore: boolean
  error: string | null
  hasMore: boolean
}>()

// Emit para enviar eventos al padre
const emit = defineEmits<{
  (e: 'refresh'): void
  (e: 'loadMore'): void
}>()

// Objeto reactivo para filtros
const filter = ref({
  ticker: '',
  rating: '',
  brokerage: '',
  limit: '',
})

// Computed para obtener ratings únicos en los stocks (para llenar dropdown)
const availableRatings = computed(() => {
  const ratingsSet = new Set(props.stocks.map((stock) => stock.ratingTo))
  return Array.from(ratingsSet).sort()
})

// Función para resetear filtros a valores vacíos
function resetFilters() {
  filter.value = {
    ticker: '',
    rating: '',
    brokerage: '',
    limit: '',
  }
}

// Computed para filtrar y limitar los stocks según filtros y límite
const filteredStocks = computed(() => {
  let filtered = props.stocks

  if (filter.value.ticker) {
    filtered = filtered.filter((stock) =>
      stock.ticker.toLowerCase().includes(filter.value.ticker.toLowerCase()),
    )
  }
  if (filter.value.rating) {
    filtered = filtered.filter((stock) => stock.ratingTo === filter.value.rating)
  }
  if (filter.value.brokerage) {
    filtered = filtered.filter((stock) =>
      stock.brokerage.toLowerCase().includes(filter.value.brokerage.toLowerCase()),
    )
  }
  if (filter.value.limit) {
    const limitNumber = parseInt(filter.value.limit)
    if (!isNaN(limitNumber) && limitNumber > 0) {
      filtered = filtered.slice(0, limitNumber)
    }
  }
  return filtered
})

// Función para emitir evento 'loadMore' para cargar más datos
function handleLoadMore() {
  emit('loadMore')
}

// Funciones para formatear y asignar clases dinámicas según valores

// Clase CSS para tipo de acción (Buy, Hold, Sell)
function getActionClass(action: string) {
  switch (action.toLowerCase()) {
    case 'buy':
      return 'bg-green-100 text-green-800'
    case 'hold':
      return 'bg-yellow-100 text-yellow-800'
    case 'sell':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

// Formatea texto de acción (capitaliza)
function formatAction(action: string) {
  if (!action) return ''
  return action.charAt(0).toUpperCase() + action.slice(1).toLowerCase()
}

// Clase CSS para rating (ejemplo, lo mismo que acción pero con otros colores)
function getRatingClass(rating: string) {
  switch (rating.toLowerCase()) {
    case 'strong buy':
      return 'bg-green-200 text-green-900'
    case 'buy':
      return 'bg-green-100 text-green-800'
    case 'hold':
      return 'bg-yellow-100 text-yellow-800'
    case 'sell':
      return 'bg-red-100 text-red-800'
    case 'strong sell':
      return 'bg-red-200 text-red-900'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

// Formatea valores monetarios, asumiendo string con número
function formatCurrency(value: string) {
  const num = parseFloat(value)
  if (isNaN(num)) return value
  return '$' + num.toFixed(2)
}

// Clase para cambio entre targetFrom y targetTo (si sube o baja)
function getTargetChangeClass(from: string, to: string) {
  const fromNum = parseFloat(from)
  const toNum = parseFloat(to)
  if (isNaN(fromNum) || isNaN(toNum)) return 'text-gray-700'

  if (toNum > fromNum) return 'text-green-600 font-semibold'
  else if (toNum < fromNum) return 'text-red-600 font-semibold'
  else return 'text-gray-700'
}
</script>
