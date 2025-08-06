<template>
  <div v-if="stock" class="space-y-6">
    <div class="flex items-center space-x-4">
      <div class="w-16 h-16 bg-gray-700 rounded-lg flex items-center justify-center">
        <span class="text-2xl font-bold text-white">{{ stock.ticker.charAt(0) }}</span>
      </div>
      <div>
        <h2 class="text-2xl font-bold text-white">{{ stock.ticker }}</h2>
        <p class="text-gray-400">{{ stock.company }}</p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="bg-gray-700 rounded-lg p-4">
        <h3 class="text-lg font-semibold text-white mb-3">Target Price</h3>
        <div class="space-y-2">
          <div class="flex justify-between">
            <span class="text-gray-400">From:</span>
            <span class="text-white font-medium">{{ stock.target_from }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-400">To:</span>
            <span class="text-white font-medium">{{ stock.target_to }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-gray-400">Change:</span>
            <div class="flex items-center">
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
      </div>

      <div class="bg-gray-700 rounded-lg p-4">
        <h3 class="text-lg font-semibold text-white mb-3">Analyst Rating</h3>
        <div class="space-y-2">
          <div class="flex justify-between">
            <span class="text-gray-400">From:</span>
            <span class="text-white">{{ stock.rating_from }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-400">To:</span>
            <span class="px-2 py-1 text-xs rounded-full" :class="getRatingColor(stock.rating_to)">
              {{ stock.rating_to }}
            </span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-400">Action:</span>
            <span class="text-white">{{ stock.action }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-400">Brokerage:</span>
            <span class="text-white">{{ stock.brokerage }}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-gray-700 rounded-lg p-4">
      <h3 class="text-lg font-semibold text-white mb-3">Additional Information</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex justify-between">
          <span class="text-gray-400">Analysis Date:</span>
          <span class="text-white">{{ formatTime(stock.time, true) }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-400">Created:</span>
          <span class="text-white">{{ formatTime(stock.created_at, true) }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-400">Updated:</span>
          <span class="text-white">{{ formatTime(stock.updated_at, true) }}</span>
        </div>
      </div>
    </div>
  </div>

  <div v-else-if="loading" class="flex justify-center py-8">
    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-green-500"></div>
  </div>

  <div v-else-if="error" class="text-center py-8">
    <p class="text-red-400">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import type { Stock } from '../../types/api'
import { getArrowPath, getRatingColor, getChangeColor, getFallbackChangePercentage, formatTime } from '../../utils/stock'

interface Props {
  stock: Stock | null
  loading: boolean
  error: string | null
}

defineProps<Props>()
</script> 