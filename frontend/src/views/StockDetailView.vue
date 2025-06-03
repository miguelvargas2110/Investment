<template>
  <div class="container mx-auto px-4 py-8">
    <!-- Encabezado de la aplicación -->
    <AppHeader />

    <!-- Enlace para volver a la lista de todas las acciones -->
    <div class="mb-6">
      <router-link to="/" class="text-blue-600 hover:text-blue-800 flex items-center">
        <!-- Icono de flecha hacia la izquierda -->
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5 mr-1"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            fill-rule="evenodd"
            d="M9.707 16.707a1 1 0 01-1.414 0l-6-6a1 1 0 010-1.414l6-6a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l4.293 4.293a1 1 0 010 1.414z"
            clip-rule="evenodd"
          />
        </svg>
        Back to all stocks
      </router-link>
    </div>

    <!-- Componente que muestra los detalles de una acción específica -->
    <StockDetail />

    <!-- Pie de página de la aplicación -->
    <AppFooter class="mt-12" />
  </div>
</template>

<script setup lang="ts">
// Importación del encabezado y pie de página reutilizables
import AppHeader from '@/components/ui/AppHeader.vue'
import AppFooter from '@/components/ui/AppFooter.vue'

// Importación del componente que muestra el detalle de una acción
import StockDetail from '@/components/stock/StockDetail.vue'
</script>
