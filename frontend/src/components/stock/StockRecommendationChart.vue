<template>
  <div class="h-96 relative">
    <!-- Aquí va el canvas donde Chart.js dibuja el gráfico -->
    <canvas ref="chartCanvas"></canvas>

    <!-- Si está cargando, mostramos un spinner con fondo semitransparente -->
    <div
      v-if="loading"
      class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-70"
    >
      <LoadingSpinner />
    </div>

    <!-- Si hay error, mostramos mensaje y botón para reintentar -->
    <div
      v-if="error"
      class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-70"
    >
      <p class="text-red-500">{{ error }}</p>
      <button
        @click="initChart"
        class="ml-2 px-3 py-1 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        Reintentar
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
// Importamos cosas de Vue para reactividad y ciclo de vida
import { ref, onMounted, watch } from 'vue'
// Importamos Chart.js y sus componentes necesarios
import { Chart, registerables } from 'chart.js'
// Importamos adaptador para que Chart.js entienda fechas
import 'chartjs-adapter-date-fns'
// Importamos nuestro store para acciones bursátiles
import { useStockStore } from '@/stores/stockStore'
// Componente spinner para mostrar mientras carga
import LoadingSpinner from '../ui/LoadingSpinner.vue'

// Registramos todos los componentes de Chart.js para usarlos
Chart.register(...registerables)

// Definimos las props que recibe el componente, aquí solo ticker (símbolo de acción)
const props = defineProps<{ ticker: string }>()

// Referencia al elemento canvas donde se dibujará el gráfico
const chartCanvas = ref<HTMLCanvasElement | null>(null)
// Guardamos la instancia del gráfico para poder destruirla y actualizarla
const chartInstance = ref<Chart<'line', { x: number; y: number }[], unknown> | null>(null)

// Variables reactivas para mostrar estado de carga y error
const loading = ref(false)
const error = ref<string | null>(null)

// Obtenemos el store para poder llamar la función que trae datos de recomendaciones
const store = useStockStore()

// Función para formatear números como dólares (USD) con dos decimales
const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(value)
}

// Función que inicializa el gráfico
const initChart = async () => {
  // Si no tenemos el canvas, no hacemos nada
  if (!chartCanvas.value) return

  // Indicamos que estamos cargando y borramos errores previos
  loading.value = true
  error.value = null

  try {
    // Pedimos al store que traiga el historial de recomendaciones para el ticker actual
    await store.fetchRecommendationHistory(props.ticker)

    // Obtenemos los datos que trajo el store
    const data = store.recommendationHistory

    // Si no hay datos, lanzamos error para mostrarlo
    if (!data || data.length === 0) {
      throw new Error('No hay datos históricos disponibles para este stock.')
    }

    // Procesamos los datos para que Chart.js pueda dibujarlos
    // Cada recomendación será una línea con dos puntos: desde precio inicial hasta precio objetivo + 90 días
    const datasets = data.map((entry, index) => {
      // Convertimos los precios a número (quitando símbolos raros si hubiera)
      const fromPrice = parseFloat(String(entry.target_from).replace(/[^\d.-]/g, ''))
      const toPrice = parseFloat(String(entry.target_to).replace(/[^\d.-]/g, ''))
      // Convertimos la fecha de la recomendación a timestamp (milisegundos)
      const date = new Date(entry.time)
      const dateTimestamp = date.getTime()
      // Calculamos fecha objetivo sumando 90 días en milisegundos
      const toTimestamp = dateTimestamp + 90 * 24 * 60 * 60 * 1000

      // Devolvemos el dataset que representa esta recomendación
      return {
        label: `Rec. ${index + 1} (${date.toLocaleDateString()})`, // Nombre de la línea
        data: [
          { x: dateTimestamp, y: fromPrice }, // Punto inicio
          { x: toTimestamp, y: toPrice }, // Punto final (objetivo)
        ],
        borderColor: index % 2 === 0 ? '#3B82F6' : '#10B981', // Alternar color azul o verde
        backgroundColor: 'transparent', // Fondo transparente
        borderWidth: 2, // Grosor línea
        borderDash: index % 2 === 0 ? [] : [5, 5], // Líneas sólidas o punteadas alternadas
        pointBackgroundColor: '#fff', // Color interior del punto
        pointBorderColor: index % 2 === 0 ? '#3B82F6' : '#10B981', // Borde del punto
        pointRadius: 5, // Tamaño del punto
        pointHoverRadius: 7, // Tamaño al pasar el mouse
        pointBorderWidth: 2, // Grosor borde del punto
        tension: 0.1, // Curvatura de la línea
        segment: {
          // Para que el segmento inicial sea sólido y el resto punteado
          borderDash: (ctx: { p0DataIndex: number }) => (ctx.p0DataIndex === 0 ? [] : [5, 5]),
        },
      }
    })

    // Si ya había un gráfico creado, lo destruimos para crear uno nuevo limpio
    if (chartInstance.value) {
      chartInstance.value.destroy()
    }

    // Obtenemos el contexto 2D del canvas para dibujar
    const ctx = chartCanvas.value.getContext('2d')
    if (!ctx) throw new Error('No se pudo obtener el contexto del canvas.')

    // Creamos la instancia del gráfico con configuración
    chartInstance.value = new Chart<'line', { x: number; y: number }[], unknown>(ctx, {
      type: 'line', // Tipo línea
      data: {
        datasets: datasets, // Datos procesados arriba
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          // Título arriba del gráfico
          title: {
            display: true,
            text: `Evolución de Recomendaciones para ${props.ticker}`,
            font: {
              size: 16,
            },
            padding: {
              top: 10,
              bottom: 20,
            },
          },
          // Tooltip personalizado para mostrar precios formateados
          tooltip: {
            callbacks: {
              label: function (context) {
                const dataset = context.dataset
                const dataIndex = context.dataIndex
                const value = context.parsed.y

                if (dataIndex === 0) {
                  return `${dataset.label} - Inicio: ${formatCurrency(value)}`
                } else {
                  return `${dataset.label} - Objetivo: ${formatCurrency(value)}`
                }
              },
            },
          },
          // Leyenda al pie del gráfico
          legend: {
            position: 'bottom',
            labels: {
              boxWidth: 12,
              padding: 20,
              usePointStyle: true,
              font: {
                size: 12,
              },
            },
          },
        },
        // Configuración de los ejes X y Y
        scales: {
          x: {
            type: 'time', // Eje X es temporal
            time: {
              unit: 'month', // Mostrar por mes
              tooltipFormat: 'dd MMM yyyy', // Formato tooltip
              displayFormats: {
                month: 'MMM yyyy', // Formato en el eje
              },
            },
            title: {
              display: true,
              text: 'Fecha de Recomendación', // Etiqueta eje X
            },
            grid: {
              display: false, // Sin líneas de grid en X
            },
          },
          y: {
            title: {
              display: true,
              text: 'Precio (USD)', // Etiqueta eje Y
            },
            ticks: {
              // Formateamos cada tick a moneda USD
              callback: function (value) {
                return formatCurrency(Number(value))
              },
            },
          },
        },
        interaction: {
          intersect: false,
          mode: 'index', // Mostrar tooltip al pasar cerca de un punto
        },
      },
    })
  } catch (err: unknown) {
    // Si algo sale mal, mostramos error en consola y en UI
    console.error('Error loading chart data:', err)
    error.value = err instanceof Error ? err.message : 'Error al cargar los datos del gráfico'
  } finally {
    // Siempre quitamos el estado de carga
    loading.value = false
  }
}

// Cuando se monte el componente, inicializamos el gráfico
onMounted(initChart)

// Si cambia la prop ticker, volvemos a cargar el gráfico
watch(
  () => props.ticker,
  (newTicker, oldTicker) => {
    if (newTicker !== oldTicker) {
      initChart()
    }
  },
)
</script>
