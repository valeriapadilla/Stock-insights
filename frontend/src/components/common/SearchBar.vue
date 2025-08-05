<template>
  <div class="search-bar relative">

    <div class="relative">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
      
      <input
        :value="modelValue"
        @input="handleInput"
        @keyup.enter="handleSearch"
        type="text"
        :placeholder="placeholder"
        class="input-field w-full pl-10 pr-4 py-2"
        :class="{ 'border-red-500': hasError }"
      />
      
      <div class="absolute inset-y-0 right-0 flex items-center">
        <button
          @click="handleSearch"
          class="px-3 py-2 text-gray-400 hover:text-gray-300 transition-colors"
          :disabled="!modelValue.trim()"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>
    </div>
    
    <p v-if="errorMessage" class="mt-1 text-sm text-red-400">
      {{ errorMessage }}
    </p>
  </div>
</template>

<script setup lang="ts">
interface Props {
  modelValue: string
  placeholder?: string
  hasError?: boolean
  errorMessage?: string
  debounceMs?: number
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Search by ticker or company name...',
  hasError: false,
  errorMessage: '',
  debounceMs: 300
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  search: [query: string]
}>()

let debounceTimer: NodeJS.Timeout | null = null

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  const value = target.value
  
  emit('update:modelValue', value)
  
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
  
  debounceTimer = setTimeout(() => {
    if (value.trim()) {
      emit('search', value.trim())
    }
  }, props.debounceMs)
}

const handleSearch = () => {
  if (props.modelValue.trim()) {
    emit('search', props.modelValue.trim())
  }
}
</script> 