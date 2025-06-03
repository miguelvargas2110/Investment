<template>
  <header class="bg-white shadow-sm mb-4">
    <div class="container mx-auto px-4 py-4 flex justify-between items-center">
      <router-link to="/" class="flex items-center">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-8 w-8 text-blue-600"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
          />
        </svg>
        <span class="ml-2 text-xl font-bold text-gray-900">StockAdvisor</span>
      </router-link>

      <nav class="flex space-x-6">
        <router-link
          to="/"
          class="text-gray-600 hover:text-blue-600 transition-colors"
          active-class="text-blue-600 font-medium"
        >
          Home
        </router-link>
        <router-link
          to="/recommendations"
          class="text-gray-600 hover:text-blue-600 transition-colors"
          active-class="text-blue-600 font-medium"
        >
          Recommendations
        </router-link>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts"></script>
