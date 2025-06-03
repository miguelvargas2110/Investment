<template>
  <!-- Si el valor 'to' es mayor que 'from', muestra una flecha hacia arriba (verde) -->
  <svg
    v-if="isIncrease"
    class="ml-1 h-4 w-4 text-green-500"
    fill="currentColor"
    viewBox="0 0 20 20"
  >
    <!-- Dibuja una flecha apuntando hacia arriba -->
    <path
      fill-rule="evenodd"
      d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"
      clip-rule="evenodd"
    />
  </svg>

  <!-- Si el valor 'to' es menor que 'from', muestra una flecha hacia abajo (roja) -->
  <svg
    v-else-if="isDecrease"
    class="ml-1 h-4 w-4 text-red-500"
    fill="currentColor"
    viewBox="0 0 20 20"
  >
    <!-- Dibuja una flecha apuntando hacia abajo -->
    <path
      fill-rule="evenodd"
      d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"
      clip-rule="evenodd"
    />
  </svg>
</template>

<script setup lang="ts">
// Define las props que recibe el componente: valores inicial y final como strings
const props = defineProps({
  from: { type: String, required: true }, // Valor original (ej. "$100.00")
  to: { type: String, required: true }, // Valor nuevo (ej. "$110.00")
})

// Elimina cualquier carácter no numérico (como "$", "%", ",") y convierte a número decimal
const fromValue = parseFloat(props.from.replace(/[^0-9.-]+/g, ''))
const toValue = parseFloat(props.to.replace(/[^0-9.-]+/g, ''))

// Verifica si el valor aumentó
const isIncrease = toValue > fromValue

// Verifica si el valor disminuyó
const isDecrease = toValue < fromValue

// Si 'fromValue' y 'toValue' son iguales, no se muestra ninguna flecha
</script>
