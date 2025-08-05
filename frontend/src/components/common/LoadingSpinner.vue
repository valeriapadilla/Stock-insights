<template>
  <div class="loading-spinner" :class="containerClass">
    <div class="spinner" :class="sizeClass">
      <div class="animate-spin rounded-full border-4 border-gray-300 border-t-green-500"></div>
    </div>
    <p v-if="showText" class="mt-2 text-gray-400 text-sm">{{ text }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  size?: 'sm' | 'md' | 'lg'
  text?: string
  showText?: boolean
  fullHeight?: boolean
  centered?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  text: 'Loading...',
  showText: true,
  fullHeight: false,
  centered: true
})

const sizeClass = computed(() => ({
  'sm': 'w-4 h-4',
  'md': 'w-8 h-8',
  'lg': 'w-12 h-12'
})[props.size])

const containerClass = computed(() => ({
  'min-h-screen': props.fullHeight,
  'py-8': !props.fullHeight,
  'flex flex-col items-center justify-center': props.centered,
  'flex justify-center': !props.centered
}))
</script>

<style scoped>
.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.spinner {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style> 