<template>
  <div class="stock-card bg-gray-800 border border-gray-700 rounded-lg p-4 hover:border-gray-600 transition-colors">
    <div class="flex items-start justify-between">

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
        
        <div class="flex flex-wrap gap-2 mb-3">
          <span class="px-2 py-1 text-xs rounded-full" :class="getRatingColor(stock.rating_to)">
            {{ stock.rating_to }}
          </span>
          <span class="px-2 py-1 text-xs rounded-full bg-gray-600 text-gray-300">
            {{ stock.brokerage }}
          </span>
        </div>
        
        <p class="text-sm text-gray-300">{{ stock.action }}</p>
      </div>

      <div class="text-right">
        <div class="mb-2">
          <p class="text-xs text-gray-400">Target Price</p>
          <p class="text-lg font-bold text-white">{{ stock.target_to }}</p>
          <p class="text-xs text-gray-500">from {{ stock.target_from }}</p>
        </div>
        
        <div class="flex items-center justify-end">
          <svg 
            class="w-4 h-4 mr-1" 
            :class="getChangeColor(stock.change_percent, stock.target_from, stock.target_to)" 
            fill="currentColor" 
            viewBox="0 0 20 20"
          >
            <path 
              fill-rule="evenodd" 
              :d="getArrowPath(stock.change_percent)" 
              clip-rule="evenodd" 
            />
          </svg>
          <span class="text-sm font-medium" :class="getChangeColor(stock.change_percent, stock.target_from, stock.target_to)">
            {{ stock.change_percent || getFallbackChangePercentage(stock.target_from, stock.target_to) }}
          </span>
        </div>
      </div>
    </div>

    <div class="mt-3 pt-3 border-t border-gray-700">
      <p class="text-xs text-gray-500">
        {{ formatTime(stock.time) }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Stock } from '../../types/api'
import { getArrowPath, getRatingColor, getChangeColor, getFallbackChangePercentage, formatTime } from '../../utils/stock'

interface Props {
  stock: Stock
}

defineProps<Props>()
</script> 