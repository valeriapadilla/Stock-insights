<template>
  <div class="filter-bar">
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
          <FilterDropdown
            v-model="selectedRating"
            :options="ratingOptions"
            placeholder="All Ratings"
          />
        </div>

        <div class="w-48">
          <FilterDropdown
            v-model="selectedSort"
            :options="sortOptions"
            placeholder="Ticker A → Z"
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
          v-if="selectedSort"
          @click="clearSort"
          class="inline-flex items-center px-2 py-1 rounded-md text-xs bg-purple-600 text-white hover:bg-purple-700"
        >
          Sort: {{ getSortLabel(selectedSort) }}
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
import { ref, computed, watch } from 'vue'
import SearchBar from '../common/SearchBar.vue'
import FilterDropdown from './FilterDropdown.vue'
import { FILTER_OPTIONS, SORT_OPTIONS } from '../../utils/constants'

interface Props {
  modelValue: {
    search: string
    rating: string
    sort_by: string
    order: 'asc' | 'desc'
  }
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [filters: Props['modelValue']]
  search: [query: string]
}>()

const searchQuery = ref(props.modelValue.search)
const selectedRating = ref(props.modelValue.rating)
const selectedSort = ref(props.modelValue.sort_by || 'ticker_asc')

const ratingOptions = FILTER_OPTIONS.RATINGS
const sortOptions = SORT_OPTIONS

const hasActiveFilters = computed(() => {
  return searchQuery.value || selectedRating.value || selectedSort.value !== 'ticker_asc'
})

watch(selectedRating, () => {
  updateFilters()
})

watch(selectedSort, () => {
  updateFilters()
})

const handleSearch = (query: string) => {
  emit('search', query)
  updateFilters()
}

const updateFilters = () => {
  emit('update:modelValue', {
    search: searchQuery.value,
    rating: selectedRating.value,
    sort_by: selectedSort.value,
    order: selectedSort.value.includes('_desc') ? 'desc' : 'asc'
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

const clearSort = () => {
  selectedSort.value = 'ticker_asc'
  updateFilters()
}

const clearAllFilters = () => {
  searchQuery.value = ''
  selectedRating.value = ''
  selectedSort.value = 'ticker_asc'
  updateFilters()
}

const getRatingLabel = (value: string) => {
  return ratingOptions.find(opt => opt.value === value)?.label || value
}

const getSortLabel = (value: string) => {
  return sortOptions.find(opt => opt.value === value)?.label || value
}
</script> 