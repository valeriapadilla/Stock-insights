<template>
  <div class="error-message bg-red-900/20 border border-red-500/50 rounded-lg p-4" :class="containerClass">
    <div class="flex items-start">
      <div class="flex-shrink-0">
        <svg class="w-5 h-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
        </svg>
      </div>
      
      <div class="ml-3 flex-1">
        <h3 v-if="title" class="text-sm font-medium text-red-200 mb-1">
          {{ title }}
        </h3>
        <p class="text-sm text-red-300">
          {{ message }}
        </p>
        
        <div v-if="showRetry" class="mt-3">
          <button 
            @click="$emit('retry')"
            class="text-sm bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded-md transition-colors"
          >
            Intentar de nuevo
          </button>
        </div>
      </div>
      
      <div v-if="showClose" class="flex-shrink-0 ml-3">
        <button 
          @click="$emit('close')"
          class="text-red-400 hover:text-red-300 transition-colors"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  message: string
  title?: string
  type?: 'error' | 'warning' | 'info'
  showRetry?: boolean
  showClose?: boolean
  fullWidth?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'error',
  showRetry: false,
  showClose: false,
  fullWidth: false
})

const emit = defineEmits<{
  retry: []
  close: []
}>()

const containerClass = computed(() => ({
  'w-full': props.fullWidth,
  'max-w-md': !props.fullWidth
}))
</script> 