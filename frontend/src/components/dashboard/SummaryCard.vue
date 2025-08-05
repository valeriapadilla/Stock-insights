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
          <component :is="icon" class="w-6 h-6" :class="iconColor" />
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
  icon?: string
  showChange?: boolean
  changeValue?: number
  changeType?: 'up' | 'down'
  changeText?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'default',
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