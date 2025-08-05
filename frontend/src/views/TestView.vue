<template>
  <div class="min-h-screen bg-gray-900 p-6">
    <div class="max-w-7xl mx-auto">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-white mb-2">StockInsights</h1>
        <p class="text-gray-400">Test Dashboard</p>
      </div>

      <!-- Summary Cards -->
      <SummaryCards />
    </div>
  </div>
</template>

<script setup lang="ts">
import SummaryCards from '../components/dashboard/SummaryCards.vue'
</script> 