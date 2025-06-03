<template>
  <div class="container mx-auto px-4 py-8">
    <!-- Cabecera -->
    <AppHeader />

    <!-- Título y descripción -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Top Stock Recommendations</h1>
      <p class="text-gray-600">
        Our algorithm analyzes multiple factors to bring you the best investment opportunities.
      </p>
    </div>

    <!-- Carga: Spinner -->
    <div v-if="loading" class="flex justify-center py-12">
      <LoadingSpinner />
    </div>

    <!-- Error -->
    <ErrorMessage v-else-if="error" :message="error" class="mb-6" />

    <!-- Lista de recomendaciones -->
    <div v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="(stock, index) in topRecommendations"
          :key="stock.ticker"
          class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300"
        >
          <div class="p-5">
            <!-- Encabezado con ticker, posición y rating -->
            <div class="flex justify-between items-start">
              <div>
                <div class="flex items-center mb-1">
                  <span class="text-xl font-bold text-gray-900 mr-2">{{ stock.ticker }}</span>
                  <span
                    class="text-xs px-2 py-1 rounded-full"
                    :class="getRatingColor(stock.rating_to)"
                  >
                    #{{ index + 1 }}
                  </span>
                </div>
                <p class="text-sm text-gray-600">{{ stock.company }}</p>
              </div>
              <span
                class="px-2 py-1 text-xs font-semibold rounded-full"
                :class="getRatingColor(stock.rating_to)"
              >
                {{ stock.rating_to }}
              </span>
            </div>

            <!-- Detalles de la recomendación -->
            <div class="mt-4">
              <div class="flex justify-between text-sm mb-2">
                <span class="text-gray-500">Target Price:</span>
                <span class="font-medium">{{ stock.target_from }} - {{ stock.target_to }}</span>
              </div>

              <div class="flex justify-between text-sm mb-2">
                <span class="text-gray-500">Brokerage:</span>
                <span class="font-medium">{{ stock.brokerage }}</span>
              </div>

              <div class="flex justify-between text-sm mb-2">
                <span class="text-gray-500">Action:</span>
                <span class="font-medium">{{ stock.action }}</span>
              </div>

              <div class="flex justify-between text-sm">
                <span class="text-gray-500">Last Updated:</span>
                <span class="font-medium">{{ formatDate(stock.time) }}</span>
              </div>
            </div>

            <!-- Enlace a detalle -->
            <div class="mt-4 pt-4 border-t border-gray-100">
              <router-link
                :to="`/stocks/${stock.ticker}`"
                class="w-full block text-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
              >
                View Details
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Pie de página -->
    <AppFooter class="mt-12" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { Chart, registerables } from 'chart.js'
import AppHeader from '@/components/ui/AppHeader.vue'
import AppFooter from '@/components/ui/AppFooter.vue'
import LoadingSpinner from '@/components/ui/LoadingSpinner.vue'
import ErrorMessage from '@/components/ui/ErrorMessage.vue'
import { useStockStore } from '@/stores/stockStore'

// Registro de componentes de Chart.js
Chart.register(...registerables)

const stockStore = useStockStore()
const loading = ref(true) // Estado de carga
const error = ref<string | null>(null) // Estado de error
const metricsChart = ref<HTMLCanvasElement | null>(null) // Referencia para un posible radar chart

// Interfaz de los datos de recomendación
interface StockRecommendation {
  ticker: string
  company: string
  rating_to: string
  target_from: number | string
  target_to: number | string
  brokerage: string
  action: string
  time: string
}

const topRecommendations = ref<StockRecommendation[]>([]) // Lista de recomendaciones

// Devuelve clases de colores según el tipo de recomendación
const getRatingColor = (rating: string) => {
  const ratingLower = rating.toLowerCase()
  if (ratingLower.includes('buy') || ratingLower.includes('outperform')) {
    return 'bg-green-100 text-green-800'
  }
  if (ratingLower.includes('sell') || ratingLower.includes('underperform')) {
    return 'bg-red-100 text-red-800'
  }
  return 'bg-yellow-100 text-yellow-800'
}

// Formatea fechas al estilo "Jun 3, 2025"
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

// Inicializa un gráfico radar con datos ficticios (puede eliminarse si no se usa)
const initMetricsChart = () => {
  if (!metricsChart.value) return
  const ctx = metricsChart.value.getContext('2d')
  if (!ctx) return

  new Chart(ctx, {
    type: 'radar',
    data: {
      labels: ['Growth', 'Value', 'Momentum', 'Risk', 'Quality', 'Sentiment'],
      datasets: [
        {
          label: 'Stock Metrics',
          data: [85, 72, 68, 45, 78, 82],
          backgroundColor: 'rgba(59, 130, 246, 0.2)',
          borderColor: 'rgba(59, 130, 246, 1)',
          borderWidth: 2,
          pointBackgroundColor: 'rgba(59, 130, 246, 1)',
          pointRadius: 4,
        },
      ],
    },
    options: {
      scales: {
        r: {
          angleLines: { display: true },
          suggestedMin: 0,
          suggestedMax: 100,
        },
      },
      plugins: {
        legend: { display: false },
      },
    },
  })
}

// Lógica para obtener las recomendaciones
const fetchRecommendations = async () => {
  try {
    loading.value = true
    error.value = null
    await stockStore.fetchTopStocks() // Llama a la store para traer datos
    topRecommendations.value = stockStore.topStocks
  } catch (err) {
    error.value = 'Failed to load recommendations. Please try again later.'
    console.error(err)
  } finally {
    loading.value = false
    await nextTick()
    initMetricsChart() // Inicializa el gráfico después del render (opcional)
  }
}

// Al montar el componente, se hace la petición
onMounted(() => {
  fetchRecommendations()
})
</script>
