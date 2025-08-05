<template>
  <div class="flex justify-center items-center" :class="containerClass">
    <div class="loading-spinner" :class="sizeClass">
      <div class="animate-spin rounded-full border-4 border-gray-300 border-t-green-500"></div>
      <p v-if="showText" class="mt-2 text-gray-400 text-sm">{{ text }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  size?: 'sm' | 'md' | 'lg'
  text?: string
  showText?: boolean
  fullHeight?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  text: 'Loading...',
  showText: true,
  fullHeight: false
})

const sizeClass = computed(() => ({
  'sm': 'w-4 h-4',
  'md': 'w-8 h-8',
  'lg': 'w-12 h-12'
})[props.size])

const containerClass = computed(() => ({
  'min-h-screen': props.fullHeight,
  'py-8': !props.fullHeight
}))
</script>

<style scoped>
.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style> 