<template>
  <!-- Contenedor principal con fondo blanco, padding, bordes redondeados y sombra -->
  <div class="bg-white p-4 rounded-lg shadow">
    <!-- Grid responsiva: 1 columna en pantallas pequeñas, 3 columnas en pantallas medianas en adelante -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Campo de entrada para el ticker -->
      <div>
        <label for="ticker" class="block text-sm font-medium text-gray-700 mb-1">Ticker</label>
        <input
          id="ticker"
          v-model="localFilter.ticker"
          type="text"
          placeholder="e.g. AAPL"
          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <!-- Campo de selección para el tipo de rating -->
      <div>
        <label for="rating" class="block text-sm font-medium text-gray-700 mb-1">Rating</label>
        <select
          id="rating"
          v-model="localFilter.rating"
          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        >
          <!-- Opción por defecto (sin filtro) -->
          <option value="">All Ratings</option>
          <option value="buy">Buy</option>
          <option value="outperform">Outperform</option>
          <option value="neutral">Neutral</option>
          <option value="underperform">Underperform</option>
          <option value="sell">Sell</option>
        </select>
      </div>

      <!-- Campo para seleccionar el número de elementos por página -->
      <div>
        <label for="limit" class="block text-sm font-medium text-gray-700 mb-1"
          >Items per page</label
        >
        <select
          id="limit"
          v-model="localFilter.limit"
          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="10">10</option>
          <option value="20">20</option>
          <option value="50">50</option>
          <option value="100">100</option>
        </select>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue' // Importa funciones reactivas de Vue

// Define las propiedades que recibe el componente desde el padre
const props = defineProps<{
  filter: {
    ticker: string // Filtro por ticker (símbolo de la acción)
    rating: string // Filtro por tipo de recomendación
    limit: number // Límite de resultados por página
  }
}>()

// Define los eventos que puede emitir el componente
const emit = defineEmits(['filter-changed']) // Se usa para notificar cambios al padre

// Crea una copia local del filtro recibido por props para modificarlo internamente
const localFilter = ref({
  ticker: props.filter.ticker,
  rating: props.filter.rating,
  limit: props.filter.limit,
})

// Observa los cambios en `localFilter` y emite el evento `filter-changed` al padre
watch(
  localFilter,
  (newValue) => {
    emit('filter-changed', newValue) // Envía el filtro actualizado al componente padre
  },
  { deep: true }, // Observa cambios profundos en objetos anidados
)
</script>
