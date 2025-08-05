<template>
  <div class="navigation-tabs mb-6">
    <div class="flex space-x-1 bg-gray-800 p-1 rounded-lg">
      <button
        @click="$emit('update:activeTab', 'stocks')"
        class="flex-1 px-6 py-3 rounded-md font-medium transition-colors"
        :class="activeTab === 'stocks' 
          ? 'bg-green-500 text-white shadow-sm' 
          : 'text-gray-400 hover:text-gray-300'"
      >
        All Stocks
      </button>
      
      <button
        @click="$emit('update:activeTab', 'recommendations')"
        class="flex-1 px-6 py-3 rounded-md font-medium transition-colors"
        :class="activeTab === 'recommendations' 
          ? 'bg-green-500 text-white shadow-sm' 
          : 'text-gray-400 hover:text-gray-300'"
      >
        Daily Recommendations
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  activeTab: 'stocks' | 'recommendations'
}

defineProps<Props>()

defineEmits<{
  'update:activeTab': [tab: 'stocks' | 'recommendations']
}>()
</script> 