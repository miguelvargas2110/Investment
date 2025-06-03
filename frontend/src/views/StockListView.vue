<template>
  <!-- Contenedor principal centrado y con padding -->
  <div class="container mx-auto px-4 py-8">
    <!-- Encabezado de la aplicación -->
    <AppHeader />

    <!-- Título y descripción de la sección de recomendaciones -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Stock Recommendations</h1>
      <p class="text-gray-600">
        Explore the latest stock recommendations from top financial analysts
      </p>
    </div>

    <!-- Filtro de búsqueda y clasificación -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <StockFilter :filter="filter" @filter-changed="handleFilterChange" />
    </div>

    <!-- Spinner de carga mientras se obtienen los datos -->
    <div v-if="loading" class="flex justify-center py-12">
      <LoadingSpinner size="lg" />
    </div>

    <!-- Sección que se muestra cuando la carga ha terminado -->
    <template v-else>
      <!-- Mensaje de error con opción de reintentar -->
      <div v-if="error" class="mb-6">
        <ErrorMessage :message="error" @retry="fetchStocks" />
      </div>

      <template v-else>
        <!-- Lista de acciones en formato grid -->
        <div
          v-if="stocks && stocks.length > 0"
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8"
        >
          <StockCard
            v-for="stock in stocks"
            :key="`${stock.ticker}-${stock.time}`"
            :stock="stock"
            @click="viewStockDetail(stock.ticker)"
          />
        </div>

        <!-- Mensaje si no se encuentran acciones con los filtros actuales -->
        <div v-else class="text-center py-12 bg-white rounded-lg shadow">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-12 w-12 mx-auto text-gray-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <h3 class="mt-4 text-lg font-medium text-gray-900">No stocks found</h3>
          <p class="mt-1 text-gray-500">Try adjusting your search or filter criteria</p>
          <button
            @click="resetFilters"
            class="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
          >
            Reset Filters
          </button>
        </div>

        <!-- Botón para cargar más acciones si hay más disponibles -->
        <div v-if="hasMore" class="flex justify-center mt-8">
          <button
            @click="loadMore"
            :disabled="loadingMore"
            class="px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <!-- Spinner cuando se están cargando más acciones -->
            <span v-if="loadingMore">
              <svg
                class="animate-spin -ml-1 mr-2 h-5 w-5 text-white inline"
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
              Loading...
            </span>
            <span v-else>Load More Stocks</span>
          </button>
        </div>
      </template>
    </template>

    <!-- Pie de página de la aplicación -->
    <AppFooter class="mt-12" />
  </div>
</template>

<script setup lang="ts">
// Composables de Vue
import { ref, onMounted, watch, computed } from 'vue'
import { useRouter } from 'vue-router'

// Store de acciones
import { useStockStore } from '@/stores/stockStore'

// Componentes reutilizables y específicos de acciones
import AppHeader from '@/components/ui/AppHeader.vue'
import AppFooter from '@/components/ui/AppFooter.vue'
import StockCard from '@/components/stock/StockCard.vue'
import StockFilter from '@/components/stock/StockFilter.vue'
import LoadingSpinner from '@/components/ui/LoadingSpinner.vue'
import ErrorMessage from '@/components/ui/ErrorMessage.vue'

// Instancia del router y store
const router = useRouter()
const stockStore = useStockStore()

// Estado reactivo del componente
const loading = ref(true)
const loadingMore = ref(false)
const error = ref<string | null>(null)

// Estado del filtro: ticker, rating, límite por página, y número de página
const filter = ref({
  ticker: '',
  rating: '',
  limit: 20,
  page: 1,
})

// Propiedades computadas
const stocks = computed(() => stockStore.stocks) // Acciones desde el store
const hasMore = computed(() => stockStore.hasMore) // Si hay más acciones para cargar

// Manejar cambios en los filtros del usuario
const handleFilterChange = (newFilter: Partial<typeof filter.value>) => {
  filter.value = { ...filter.value, ...newFilter, page: 1 }
}

// Reiniciar filtros a sus valores por defecto
const resetFilters = () => {
  filter.value = {
    ticker: '',
    rating: '',
    limit: 20,
    page: 1,
  }
}

// Navegar al detalle de una acción
const viewStockDetail = (ticker: string) => {
  router.push(`/stocks/${ticker}`)
}

// Obtener acciones desde el store
const fetchStocks = async () => {
  try {
    loading.value = true
    error.value = null
    await stockStore.fetchStocks(filter.value)
  } catch (err) {
    error.value = 'Failed to load stocks. Please try again later.'
    console.error('Error fetching stocks:', err)
  } finally {
    loading.value = false
  }
}

// Cargar más acciones al presionar "Load More"
const loadMore = async () => {
  try {
    loadingMore.value = true
    filter.value.page += 1
    await stockStore.loadMoreStocks()
  } catch (err) {
    error.value = 'Failed to load more stocks.'
    console.error('Error loading more stocks:', err)
  } finally {
    loadingMore.value = false
  }
}

// Ejecutar fetchStocks cuando los filtros cambian
watch(
  filter,
  () => {
    fetchStocks()
  },
  { deep: true },
)

// Ejecutar fetchStocks al montar el componente
onMounted(() => {
  fetchStocks()
})
</script>

<style scoped>
/* Transiciones para mostrar/ocultar contenido con suavidad */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Efecto de elevación en tarjetas de acciones al hacer hover */
.stock-card {
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.stock-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}
</style>
