<template>
  <div class="stock-card bg-gray-800 border border-gray-700 rounded-lg p-4 hover:border-gray-600 transition-colors">
    <div class="flex items-start justify-between">
      <!-- Left side: Ticker and company info -->
      <div class="flex-1">
        <div class="flex items-center mb-2">
          <div class="w-10 h-10 bg-gray-700 rounded-lg flex items-center justify-center mr-3">
            <span class="text-lg font-bold text-white">{{ stock.ticker.charAt(0) }}</span>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-white">{{ stock.ticker }}</h3>
            <p class="text-sm text-gray-400">{{ stock.company }}</p>
          </div>
        </div>
        
        <!-- Tags -->
        <div class="flex flex-wrap gap-2 mb-3">
          <span class="px-2 py-1 text-xs rounded-full" :class="getRatingColor(stock.rating_to)">
            {{ stock.rating_to }}
          </span>
          <span class="px-2 py-1 text-xs rounded-full bg-gray-600 text-gray-300">
            {{ stock.brokerage }}
          </span>
        </div>
        
        <!-- Action -->
        <p class="text-sm text-gray-300">{{ stock.action }}</p>
      </div>

      <!-- Right side: Price info -->
      <div class="text-right">
        <div class="mb-2">
          <p class="text-xs text-gray-400">Target Price</p>
          <p class="text-lg font-bold text-white">{{ stock.target_to }}</p>
          <p class="text-xs text-gray-500">from {{ stock.target_from }}</p>
        </div>
        
        <!-- Change indicator -->
        <div class="flex items-center justify-end">
          <svg 
            class="w-4 h-4 mr-1" 
            :class="getChangeColor()" 
            fill="currentColor" 
            viewBox="0 0 20 20"
          >
            <path fill-rule="evenodd" d="M3.293 9.707a1 1 0 010-1.414l6-6a1 1 0 011.414 0l6 6a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L4.707 9.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
          </svg>
          <span class="text-sm font-medium" :class="getChangeColor()">
            {{ getChangePercentage() }}
          </span>
        </div>
      </div>
    </div>

    <!-- Time info -->
    <div class="mt-3 pt-3 border-t border-gray-700">
      <p class="text-xs text-gray-500">
        {{ formatTime(stock.time) }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Stock } from '../../types/api'

interface Props {
  stock: Stock
}

const props = defineProps<Props>()

// Computed
const getRatingColor = (rating: string) => {
  switch (rating.toLowerCase()) {
    case 'buy':
      return 'bg-green-600 text-white'
    case 'sell':
      return 'bg-red-600 text-white'
    case 'hold':
      return 'bg-yellow-600 text-white'
    case 'neutral':
      return 'bg-blue-600 text-white'
    default:
      return 'bg-gray-600 text-gray-300'
  }
}

const getChangeColor = () => {
  // Simple logic: if target_to > target_from, it's positive
  const fromPrice = parseFloat(props.stock.target_from.replace('$', ''))
  const toPrice = parseFloat(props.stock.target_to.replace('$', ''))
  return toPrice > fromPrice ? 'text-green-400' : 'text-red-400'
}

const getChangePercentage = () => {
  const fromPrice = parseFloat(props.stock.target_from.replace('$', ''))
  const toPrice = parseFloat(props.stock.target_to.replace('$', ''))
  const change = ((toPrice - fromPrice) / fromPrice) * 100
  return `${change > 0 ? '+' : ''}${change.toFixed(1)}%`
}

const formatTime = (timeString: string) => {
  const date = new Date(timeString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script> 