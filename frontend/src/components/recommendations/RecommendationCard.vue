<template>
  <div class="recommendation-card bg-gray-800 border border-gray-700 rounded-lg p-6 hover:border-gray-600 transition-colors">
    <div class="flex items-start justify-between">
      <div class="flex-1">
        <div class="flex items-center mb-3">
          <div class="w-12 h-12 bg-gradient-to-r from-green-500 to-green-600 rounded-lg flex items-center justify-center mr-4">
            <span class="text-lg font-bold text-white">#{{ recommendation.rank }}</span>
          </div>
          <div>
            <h3 class="text-xl font-bold text-white">{{ recommendation.ticker }}</h3>
            <div class="flex items-center mt-1">
              <span class="px-2 py-1 text-xs rounded-full bg-green-600 text-white font-medium">
                Score: {{ recommendation.score }}/100
              </span>
            </div>
          </div>
        </div>
        
        <!-- Explanation -->
        <p class="text-sm text-gray-300 leading-relaxed">
          {{ recommendation.explanation }}
        </p>
      </div>

      <!-- Right side: Score -->
      <div class="text-right ml-6">
        <div class="mb-2">
          <p class="text-xs text-gray-400">AI Score</p>
          <p class="text-3xl font-bold text-green-400">{{ recommendation.score }}</p>
        </div>
        
        <!-- View Details button -->
        <button 
          @click="$emit('viewDetails', recommendation)"
          class="btn-primary text-sm px-4 py-2"
        >
          View Details
        </button>
      </div>
    </div>

    <!-- Run time info -->
    <div class="mt-4 pt-4 border-t border-gray-700">
      <p class="text-xs text-gray-500">
        Calculated: {{ formatTime(recommendation.run_at) }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Recommendation } from '../../types/api'

interface Props {
  recommendation: Recommendation
}

defineProps<Props>()

defineEmits<{
  viewDetails: [recommendation: Recommendation]
}>()

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