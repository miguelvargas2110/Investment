<template>
  <div class="container mx-auto px-4 py-8">
    <AppHeader />
    <!-- Encabezado de la app -->

    <!-- Carrusel con las mejores recomendaciones -->
    <div class="mb-8">
      <h2 class="text-2xl font-semibold text-gray-900 mb-4">Top Recommendations</h2>
      <div class="overflow-x-auto scrollbar-hide">
        <div class="flex space-x-4">
          <!-- Recorre las acciones top para mostrar cada una como un enlace -->
          <router-link
            v-for="stock in topStocks"
            :key="stock.ticker"
            :to="`/stocks/${stock.ticker}`"
            class="min-w-[200px] flex-shrink-0 bg-white rounded-xl shadow-md p-4 hover:shadow-lg transition duration-300 border"
          >
            <div class="flex flex-col justify-between h-full">
              <div>
                <h3 class="text-lg font-medium text-gray-900">{{ stock.ticker }}</h3>
                <!-- Símbolo -->
                <p class="text-sm text-gray-600 truncate">{{ stock.company }}</p>
                <!-- Nombre empresa -->
              </div>
              <!-- Muestra la calificación con colores según valor -->
              <span
                :class="getRatingClass(stock.rating_to)"
                class="mt-4 self-start px-2 py-1 text-xs font-semibold rounded-full"
              >
                {{ stock.rating_to }}
              </span>
            </div>
          </router-link>
        </div>
      </div>
    </div>

    <!-- Contenido principal con título y descripción -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Stock Recommendations</h1>
      <p class="text-gray-600">Discover the latest stock recommendations from top brokerages.</p>
    </div>

    <!-- Lista de acciones recomendadas dentro de un card -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <div class="lg:col-span-2">
        <!-- Componente para listar acciones -->
        <StockList
          :stocks="formattedStocks"
          :loading="loadingState"
          :loadingMore="false"
          :error="errorMessage"
          :hasMore="hasMoreStocks"
          @refresh="handleRefresh"
          @load-more="handleLoadMore"
          class="h-full"
        />
      </div>
    </div>

    <AppFooter class="mt-12" />
    <!-- Pie de página -->
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useStockStore } from '@/stores/stockStore' // Store para manejar datos de acciones
import AppHeader from '@/components/ui/AppHeader.vue'
import AppFooter from '@/components/ui/AppFooter.vue'
import StockList from '@/components/stock/StockList.vue'

const stockStore = useStockStore() // Instancia del store
const currentPage = ref(1) // Página actual para paginación
const hasMoreStocks = ref(true) // Indica si hay más datos para cargar

// Computed para simplificar reactividad desde el store
const loadingState = computed(() => stockStore.loading) // Estado de carga
const errorMessage = computed(() => stockStore.error ?? null) // Mensajes de error
const topStocks = computed(() => stockStore.topStocks) // Las mejores recomendaciones

// Formatea stocks para que las propiedades estén con nombres más amigables en el componente hijo
const formattedStocks = computed(() => {
  return stockStore.stocks.map((stock) => ({
    ...stock,
    ratingTo: stock.rating_to,
    targetFrom: stock.target_from.toString(),
    targetTo: stock.target_to.toString(),
  }))
})

// Refrescar lista, reseteando página y cargando datos
const handleRefresh = async () => {
  currentPage.value = 1
  hasMoreStocks.value = true
  await stockStore.fetchStocks({
    ticker: '',
    rating: '',
    limit: 20,
    page: currentPage.value,
  })
}

// Cargar más datos para paginación
const handleLoadMore = async () => {
  currentPage.value++
  try {
    await stockStore.loadMoreStocks()
    // Si la cantidad de stocks es múltiplo del límite, asumimos que hay más páginas
    hasMoreStocks.value = stockStore.stocks.length % 20 === 0
  } catch (error) {
    currentPage.value-- // Si falla, revertimos la página
    throw error
  }
}

// Función para asignar clases de color según la calificación
const getRatingClass = (rating: string) => {
  if (rating.includes('Buy') || rating.includes('Outperform')) return 'bg-green-100 text-green-800'
  if (rating.includes('Sell') || rating.includes('Underperform')) return 'bg-red-100 text-red-800'
  return 'bg-yellow-100 text-yellow-800'
}

// Al montar el componente, se hace la carga inicial de datos
onMounted(() => {
  handleRefresh()
  stockStore.fetchTopStocks()
  stockStore.fetchTickers()
})
</script>
