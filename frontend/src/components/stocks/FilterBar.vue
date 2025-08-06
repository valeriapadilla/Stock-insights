<template>
  <div class="filter-bar bg-gray-800 border border-gray-700 rounded-lg p-4">
    <div class="flex flex-col lg:flex-row gap-4">
      <div class="flex-1">
        <SearchBar
          v-model="searchQuery"
          placeholder="Search by ticker or company name..."
          @search="handleSearch"
        />
      </div>

      <div class="flex gap-3">
        <!-- Rating Filter -->
        <div class="w-32">
          <label class="block text-sm font-medium text-gray-300 mb-1">Rating</label>
          <FilterDropdown
            v-model="selectedRating"
            :options="ratingOptions"
            placeholder="All Ratings"
          />
        </div>

        <!-- Action Filter -->
        <div class="w-32">
          <label class="block text-sm font-medium text-gray-300 mb-1">Action</label>
          <FilterDropdown
            v-model="selectedAction"
            :options="actionOptions"
            placeholder="All Actions"
          />
        </div>

        <div class="w-32">
          <label class="block text-sm font-medium text-gray-300 mb-1">Price</label>
          <FilterDropdown
            v-model="selectedPriceRange"
            :options="priceRangeOptions"
            placeholder="All Prices"
          />
        </div>
      </div>
    </div>

    <div v-if="hasActiveFilters" class="mt-4 pt-4 border-t border-gray-700">
      <div class="flex flex-wrap gap-2">
        <span class="text-sm text-gray-400">Active filters:</span>
        
        <button
          v-if="searchQuery"
          @click="clearSearch"
          class="inline-flex items-center px-2 py-1 rounded-md text-xs bg-blue-600 text-white hover:bg-blue-700"
        >
          Search: "{{ searchQuery }}"
          <svg class="w-3 h-3 ml-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>

        <button
          v-if="selectedRating"
          @click="clearRating"
          class="inline-flex items-center px-2 py-1 rounded-md text-xs bg-green-600 text-white hover:bg-green-700"
        >
          Rating: {{ getRatingLabel(selectedRating) }}
          <svg class="w-3 h-3 ml-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>

        <button
          v-if="selectedAction"
          @click="clearAction"
          class="inline-flex items-center px-2 py-1 rounded-md text-xs bg-purple-600 text-white hover:bg-purple-700"
        >
          Action: {{ getActionLabel(selectedAction) }}
          <svg class="w-3 h-3 ml-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>

        <button
          v-if="selectedPriceRange"
          @click="clearPriceRange"
          class="inline-flex items-center px-2 py-1 rounded-md text-xs bg-orange-600 text-white hover:bg-orange-700"
        >
          Price: {{ getPriceRangeLabel(selectedPriceRange) }}
          <svg class="w-3 h-3 ml-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>

        <button
          @click="clearAllFilters"
          class="inline-flex items-center px-2 py-1 rounded-md text-xs bg-gray-600 text-white hover:bg-gray-700"
        >
          Clear all
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import SearchBar from '../common/SearchBar.vue'
import FilterDropdown from './FilterDropdown.vue'
import { FILTER_OPTIONS } from '../../utils/constants'

interface Props {
  modelValue: {
    search: string
    rating: string
    action: string
    priceRange: string
  }
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [filters: Props['modelValue']]
  search: [query: string]
}>()

const searchQuery = ref(props.modelValue.search)
const selectedRating = ref(props.modelValue.rating)
const selectedAction = ref(props.modelValue.action)
const selectedPriceRange = ref(props.modelValue.priceRange)

const ratingOptions = FILTER_OPTIONS.RATINGS
const actionOptions = FILTER_OPTIONS.ACTIONS
const priceRangeOptions = FILTER_OPTIONS.PRICE_RANGES

const hasActiveFilters = computed(() => {
  return searchQuery.value || selectedRating.value || selectedAction.value || selectedPriceRange.value
})

const handleSearch = (query: string) => {
  emit('search', query)
  updateFilters()
}

const updateFilters = () => {
  emit('update:modelValue', {
    search: searchQuery.value,
    rating: selectedRating.value,
    action: selectedAction.value,
    priceRange: selectedPriceRange.value
  })
}

const clearSearch = () => {
  searchQuery.value = ''
  updateFilters()
}

const clearRating = () => {
  selectedRating.value = ''
  updateFilters()
}

const clearAction = () => {
  selectedAction.value = ''
  updateFilters()
}

const clearPriceRange = () => {
  selectedPriceRange.value = ''
  updateFilters()
}

const clearAllFilters = () => {
  searchQuery.value = ''
  selectedRating.value = ''
  selectedAction.value = ''
  selectedPriceRange.value = ''
  updateFilters()
}

const getRatingLabel = (value: string) => {
  return ratingOptions.find(opt => opt.value === value)?.label || value
}

const getActionLabel = (value: string) => {
  return actionOptions.find(opt => opt.value === value)?.label || value
}

const getPriceRangeLabel = (value: string) => {
  return priceRangeOptions.find(opt => opt.value === value)?.label || value
}
</script> 