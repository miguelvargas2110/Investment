<template>
  <!-- Estado de carga: muestra un spinner -->
  <div v-if="loading" class="flex justify-center py-8">
    <LoadingSpinner />
  </div>

  <!-- Estado de error: muestra mensaje y botón para reintentar -->
  <div v-else-if="error" class="text-center py-8">
    <ErrorMessage :message="error" />
    <button
      @click="fetchStock"
      class="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
    >
      Retry
    </button>
  </div>

  <!-- Estado de éxito: muestra los detalles de la acción -->
  <div v-else-if="stock" class="bg-white rounded-lg shadow-md overflow-hidden">
    <div class="p-6">
      <!-- Encabezado con el ticker, empresa y rating -->
      <div class="flex justify-between items-start">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">{{ stock.ticker }} - {{ stock.company }}</h1>
          <p class="text-gray-600">{{ stock.brokerage }}</p>
        </div>
        <!-- Badge con color según el rating -->
        <span :class="ratingClass" class="px-3 py-1 text-sm font-semibold rounded-full">
          {{ stock.rating_to }}
        </span>
      </div>

      <!-- Contenido dividido en dos columnas -->
      <div class="mt-6 grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Columna izquierda: detalles de la recomendación -->
        <div>
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Recommendation Details</h2>
          <div class="space-y-3">
            <div class="flex justify-between">
              <span class="text-gray-600">Action:</span>
              <span class="font-medium">{{ stock.action }}</span>
            </div>

            <div class="flex justify-between">
              <span class="text-gray-600">Previous Rating:</span>
              <span class="font-medium">{{ stock.rating_from }}</span>
            </div>

            <div class="flex justify-between">
              <span class="text-gray-600">Current Rating:</span>
              <span class="font-medium">{{ stock.rating_to }}</span>
            </div>

            <div class="flex justify-between">
              <span class="text-gray-600">Price Target:</span>
              <span class="font-medium">{{ stock.target_from }} - {{ stock.target_to }}</span>
            </div>

            <div class="flex justify-between">
              <span class="text-gray-600">Date:</span>
              <span class="font-medium">{{ formatDate(stock.time) }}</span>
            </div>
          </div>
        </div>

        <!-- Columna derecha: gráfico de rendimiento -->
        <div>
          <h2 class="text-lg font-semibold text-gray-900 mb-4">Performance Chart</h2>
          <StockRecommendationChart :ticker="stock.ticker" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Importaciones necesarias
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import LoadingSpinner from '../ui/LoadingSpinner.vue'
import ErrorMessage from '../ui/ErrorMessage.vue'
import StockRecommendationChart from './StockRecommendationChart.vue'
import { useStockStore } from '@/stores/stockStore'

// Acceso a la ruta y al store de acciones
const route = useRoute()
const stockStore = useStockStore()

// Estados de carga y error
const loading = ref(true)
const error = ref<string | null>(null)

// Computed con los datos actuales del stock desde el store
const stock = computed(() => stockStore.currentStock)

// Clase dinámica basada en el tipo de recomendación
const ratingClass = computed(() => {
  const rating = stock.value?.rating_to.toLowerCase() || ''
  if (rating.includes('buy') || rating.includes('outperform')) {
    return 'bg-green-100 text-green-800'
  }
  if (rating.includes('sell') || rating.includes('underperform')) {
    return 'bg-red-100 text-red-800'
  }
  return 'bg-yellow-100 text-yellow-800'
})

// Formatear la fecha de la recomendación
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

// Función para cargar detalles del stock desde el store
const fetchStock = async () => {
  const ticker = route.params.ticker as string
  if (!ticker) {
    error.value = 'No stock ticker provided.'
    loading.value = false
    return
  }

  try {
    loading.value = true
    error.value = null
    await stockStore.fetchStockDetail(ticker)
  } catch (err) {
    error.value = 'Failed to load stock details. Please try again later.'
    console.error(err)
  } finally {
    loading.value = false
  }
}

// Carga inicial del stock cuando el componente se monta
onMounted(() => {
  fetchStock()
})

// Vuelve a cargar el stock si cambia el parámetro de la ruta
watch(
  () => route.params.ticker,
  () => {
    fetchStock()
  },
)
</script>
