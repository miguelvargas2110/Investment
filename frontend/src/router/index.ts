// Importamos funciones necesarias para crear el router con historial web
import { createRouter, createWebHistory } from 'vue-router'

// Importamos las vistas (componentes) que usaremos en las rutas
import HomeView from '@/views/HomeView.vue'
import StockDetailView from '@/views/StockDetailView.vue'
import RecommendationsView from '@/views/RecommendationsView.vue'
import StockListView from '@/views/StockListView.vue'

// Creamos la instancia del router con configuración
const router = createRouter({
  // Usamos historial de navegador para URLs limpias (sin #)
  history: createWebHistory(import.meta.env.BASE_URL),

  // Definimos las rutas disponibles en la app
  routes: [
    {
      path: '/', // Ruta principal o "home"
      name: 'home', // Nombre para referenciar esta ruta en código
      component: HomeView, // Componente que se muestra en esta ruta
    },
    {
      path: '/stocks/:ticker', // Ruta dinámica para detalle de acción, ":ticker" es parámetro
      name: 'stock-detail', // Nombre para identificar esta ruta
      component: StockDetailView, // Componente que muestra el detalle de una acción
      props: true, // Para que el parámetro "ticker" se pase como prop al componente
    },
    {
      path: '/recommendations', // Ruta para ver recomendaciones de acciones
      name: 'recommendations',
      component: RecommendationsView,
    },
    {
      path: '/stocks', // Ruta para ver la lista de todas las acciones
      name: 'stocks',
      component: StockListView,
    },
  ],
})

// Exportamos el router para usarlo en la app principal (main.ts)
export default router
