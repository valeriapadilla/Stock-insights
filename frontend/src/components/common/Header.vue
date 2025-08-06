<template>
  <header class="bg-gray-900 border-b border-gray-800">
    <div class="max-w-7xl mx-auto px-6 py-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="w-8 h-8 bg-green-500 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-white" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M3.293 9.707a1 1 0 010-1.414l6-6a1 1 0 011.414 0l6 6a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L4.707 9.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-white">StockInsights</h1>
        </div>

        <!-- Last Update Indicator -->
        <div class="flex items-center space-x-2">
          <div class="w-2 h-2 bg-green-500 rounded-full"></div>
          <span class="text-sm text-green-400 font-medium">
            {{ lastUpdateText }}
          </span>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useStocksStore } from '../../stores/stocks'

const stocksStore = useStocksStore()
const lastUpdateText = ref('Loading...')

const loadLastUpdate = async () => {
  try {
    await stocksStore.loadStocks({ limit: 1, sort_by: 'created_at', order: 'desc' })
    
    if (stocksStore.stocks.length > 0) {
      const latestStock = stocksStore.stocks[0]
      const lastUpdate = new Date(latestStock.created_at)
      
      const now = new Date()
      const diffInHours = Math.floor((now.getTime() - lastUpdate.getTime()) / (1000 * 60 * 60))
      
      if (diffInHours < 1) {
        lastUpdateText.value = 'Updated just now'
      } else if (diffInHours < 24) {
        lastUpdateText.value = `Updated ${diffInHours}h ago`
      } else {
        const diffInDays = Math.floor(diffInHours / 24)
        lastUpdateText.value = `Updated ${diffInDays}d ago`
      }
    } else {
      lastUpdateText.value = 'No data available'
    }
  } catch (error) {
    lastUpdateText.value = 'Update time unavailable'
  }
}

onMounted(() => {
  loadLastUpdate()
})
</script> 