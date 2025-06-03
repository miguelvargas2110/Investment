<template>
  <!-- Tarjeta principal con estilo visual -->
  <div
    class="bg-white rounded-xl shadow-sm overflow-hidden border border-gray-100 hover:border-blue-200 transition-all"
  >
    <div class="p-5">
      <!-- Encabezado: muestra el ticker y el nombre de la empresa -->
      <div class="flex justify-between items-start">
        <div>
          <h3 class="text-xl font-bold text-gray-900">{{ stock.ticker }}</h3>
          <!-- Símbolo bursátil -->
          <p class="text-sm text-gray-500">{{ stock.company }}</p>
          <!-- Nombre de la empresa -->
        </div>
        <!-- Etiqueta de rating con color dinámico -->
        <span
          :class="ratingClass"
          class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
        >
          {{ stock.rating_to }}
          <!-- Valor actual de la recomendación -->
        </span>
      </div>

      <!-- Detalles adicionales de la recomendación -->
      <div class="mt-4 space-y-2">
        <!-- Nombre del corredor o institución -->
        <div class="flex justify-between">
          <span class="text-sm text-gray-500">Brokerage</span>
          <span class="text-sm font-medium text-gray-900">{{ stock.brokerage }}</span>
        </div>

        <!-- Acción sugerida (Buy, Hold, Sell, etc.) -->
        <div class="flex justify-between">
          <span class="text-sm text-gray-500">Action</span>
          <span class="text-sm font-medium text-gray-900">{{ stock.action }}</span>
        </div>

        <!-- Precio objetivo, con transición visual -->
        <div class="flex justify-between">
          <span class="text-sm text-gray-500">Target Price</span>
          <span class="text-sm font-medium text-blue-600">
            {{ stock.target_from }} → {{ stock.target_to }}
          </span>
        </div>
      </div>

      <!-- Pie de tarjeta: fecha y botón -->
      <div class="mt-4 pt-4 border-t border-gray-100 flex justify-between items-center">
        <!-- Fecha formateada -->
        <span class="text-xs text-gray-500">
          {{ formatDate(stock.time) }}
        </span>
        <!-- Botón que podría abrir detalles (todavía no implementado) -->
        <button class="text-xs font-medium text-blue-600 hover:text-blue-800">View Details</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Importaciones necesarias
import { computed } from 'vue'
import type { StockRecommendation } from '@/types'

// Props que recibe el componente: un objeto de tipo StockRecommendation
const props = defineProps<{
  stock: StockRecommendation
}>()

// Computed property para asignar una clase de color al rating
const ratingClass = computed(() => {
  const rating = props.stock.rating_to.toLowerCase()
  // Recomendaciones positivas
  if (rating.includes('buy') || rating.includes('outperform')) {
    return 'bg-green-100 text-green-800'
  }
  // Recomendaciones negativas
  if (rating.includes('sell') || rating.includes('underperform')) {
    return 'bg-red-100 text-red-800'
  }
  // Recomendaciones neutras (hold, neutral, etc.)
  return 'bg-yellow-100 text-yellow-800'
})

// Función para formatear la fecha en formato corto en inglés
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>
