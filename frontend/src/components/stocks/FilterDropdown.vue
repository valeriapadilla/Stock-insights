<template>
  <div class="relative">

    <button
      @click="toggleDropdown"
      class="dropdown flex items-center justify-between w-full px-3 py-2 text-left"
      :class="{ 'border-green-500': isOpen }"
    >
      <span class="text-gray-300">{{ selectedLabel }}</span>
      <svg 
        class="w-4 h-4 text-gray-400 transition-transform"
        :class="{ 'rotate-180': isOpen }"
        fill="none" 
        stroke="currentColor" 
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>

    <!-- Dropdown menu -->
    <div 
      v-if="isOpen"
      class="absolute z-10 mt-1 w-full bg-gray-700 border border-gray-600 rounded-lg shadow-lg"
    >
      <div class="py-1">
        <button
          v-for="option in options"
          :key="option.value"
          @click="selectOption(option)"
          class="w-full px-3 py-2 text-left text-gray-300 hover:bg-gray-600 transition-colors flex items-center justify-between"
          :class="{ 'bg-gray-600': option.value === modelValue }"
        >
          <span>{{ option.label }}</span>
          <svg 
            v-if="option.value === modelValue"
            class="w-4 h-4 text-green-400" 
            fill="currentColor" 
            viewBox="0 0 20 20"
          >
            <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface FilterOption {
  value: string
  label: string
}

interface Props {
  modelValue: string
  options: FilterOption[]
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: 'Select option'
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const isOpen = ref(false)

const selectedLabel = computed(() => {
  const selected = props.options.find(opt => opt.value === props.modelValue)
  return selected?.label || props.placeholder
})

// Methods
const toggleDropdown = () => {
  isOpen.value = !isOpen.value
}

const selectOption = (option: FilterOption) => {
  emit('update:modelValue', option.value)
  isOpen.value = false
}

const closeDropdown = (event: Event) => {
  const target = event.target as Element
  if (!target.closest('.relative')) {
    isOpen.value = false
  }
}

// Lifecycle
onMounted(() => {
  document.addEventListener('click', closeDropdown)
})

onUnmounted(() => {
  document.removeEventListener('click', closeDropdown)
})
</script> 