<template>
  <div class="summary-card bg-gray-800 border border-gray-700 rounded-lg p-6">
    <div class="flex items-center justify-between">
      <div class="flex-1">
        <h3 class="text-sm font-medium text-gray-400 mb-1">{{ title }}</h3>
        <p class="text-3xl font-bold" :class="valueColor">{{ formattedValue }}</p>
        <p class="text-sm text-gray-500 mt-1">{{ description }}</p>
      </div>

      <div class="flex-shrink-0 ml-4">
        <div class="w-12 h-12 rounded-lg flex items-center justify-center" :class="iconBgColor">
          <!-- Stocks Icon -->
          <svg v-if="iconType === 'stocks'" class="w-6 h-6" :class="iconColor" fill="currentColor" viewBox="0 0 24 24">
            <path d="M3 3h18v18H3V3zm16 16V5H5v14h14z"/>
            <path d="M7 7h4v4H7V7zm6 0h4v4h-4V7zM7 13h4v4H7v-4zm6 0h4v4h-4v-4z"/>
          </svg>
          
          <!-- Recommendations Icon -->
          <svg v-else-if="iconType === 'recommendations'" class="w-6 h-6" :class="iconColor" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
          </svg>
          
          <!-- Default Icon -->
          <svg v-else class="w-6 h-6" :class="iconColor" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
          </svg>
        </div>
      </div>
    </div>

    <div v-if="showChange" class="mt-4 pt-4 border-t border-gray-700">
      <div class="flex items-center">
        <svg 
          v-if="changeIcon" 
          class="w-4 h-4 mr-1" 
          :class="changeIconColor" 
          fill="currentColor" 
          viewBox="0 0 20 20"
        >
          <path v-if="changeType === 'up'" fill-rule="evenodd" d="M3.293 9.707a1 1 0 010-1.414l6-6a1 1 0 011.414 0l6 6a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L4.707 9.707a1 1 0 01-1.414 0z" clip-rule="evenodd" />
          <path v-else fill-rule="evenodd" d="M16.707 10.293a1 1 0 010 1.414l-6 6a1 1 0 01-1.414 0l-6-6a1 1 0 111.414-1.414L9 14.586V3a1 1 0 012 0v11.586l4.293-4.293a1 1 0 011.414 0z" clip-rule="evenodd" />
        </svg>
        <span class="text-sm font-medium" :class="changeTextColor">
          {{ changeText }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  title: string
  value: number
  description: string
  type?: 'default' | 'success' | 'warning' | 'error'
  iconType?: 'stocks' | 'recommendations' | 'default'
  showChange?: boolean
  changeValue?: number
  changeType?: 'up' | 'down'
  changeText?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'default',
  iconType: 'default',
  showChange: false,
  changeType: 'up',
  changeText: ''
})

const formattedValue = computed(() => {
  if (props.value >= 1000) {
    return props.value.toLocaleString()
  }
  return props.value.toString()
})

const valueColor = computed(() => {
  switch (props.type) {
    case 'success': return 'text-green-500'
    case 'warning': return 'text-yellow-500'
    case 'error': return 'text-red-500'
    default: return 'text-white'
  }
})

const iconBgColor = computed(() => {
  switch (props.type) {
    case 'success': return 'bg-green-500/20'
    case 'warning': return 'bg-yellow-500/20'
    case 'error': return 'bg-red-500/20'
    default: return 'bg-gray-600/20'
  }
})

const iconColor = computed(() => {
  switch (props.type) {
    case 'success': return 'text-green-400'
    case 'warning': return 'text-yellow-400'
    case 'error': return 'text-red-400'
    default: return 'text-gray-400'
  }
})

const changeIcon = computed(() => props.showChange && props.changeType)

const changeIconColor = computed(() => {
  return props.changeType === 'up' ? 'text-green-400' : 'text-red-400'
})

const changeTextColor = computed(() => {
  return props.changeType === 'up' ? 'text-green-400' : 'text-red-400'
})
</script> 