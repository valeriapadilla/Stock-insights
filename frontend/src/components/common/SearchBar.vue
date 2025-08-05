<template>
  <div class="search-bar relative">
    <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
      <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
    </div>
    <input
      :value="modelValue"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      @keyup.enter="$emit('search', modelValue)"
      type="text"
      class="block w-full pl-10 pr-3 py-2 border border-gray-600 rounded-md leading-5 bg-gray-700 text-white placeholder-gray-400 focus:outline-none focus:placeholder-gray-500 focus:ring-1 focus:ring-green-500 focus:border-green-500"
      :placeholder="placeholder"
    />
  </div>
</template>

<script setup lang="ts">
interface Props {
  modelValue: string
  placeholder?: string
}

defineProps<Props>()

defineEmits<{
  'update:modelValue': [value: string]
  search: [query: string]
}>()
</script> 