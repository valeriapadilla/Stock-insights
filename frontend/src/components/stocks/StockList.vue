<template>
  <div class="stock-list">
    <div v-if="isLoading" class="space-y-4">
      <div v-for="i in 5" :key="i" class="bg-gray-800 border border-gray-700 rounded-lg p-4 animate-pulse">
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <div class="flex items-center mb-2">
              <div class="w-10 h-10 bg-gray-700 rounded-lg mr-3"></div>
              <div>
                <div class="h-5 bg-gray-700 rounded w-16 mb-1"></div>
                <div class="h-3 bg-gray-700 rounded w-24"></div>
              </div>
            </div>
            <div class="flex gap-2 mb-3">
              <div class="h-5 bg-gray-700 rounded w-12"></div>
              <div class="h-5 bg-gray-700 rounded w-20"></div>
            </div>
            <div class="h-3 bg-gray-700 rounded w-32"></div>
          </div>
          <div class="text-right">
            <div class="h-4 bg-gray-700 rounded w-16 mb-1"></div>
            <div class="h-6 bg-gray-700 rounded w-20"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Error state -->
    <ErrorMessage 
      v-else-if="error"
      :message="error"
      :showRetry="true"
      @retry="loadStocks"
    />

    <!-- Empty state -->
    <div v-else-if="stocks.length === 0" class="text-center py-12">
      <div class="text-gray-400 mb-4">
        <svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-300 mb-2">No stocks found</h3>
        <p class="text-gray-500">Try adjusting your filters or search terms.</p>
      </div>
    </div>

    <!-- Data state -->
    <div v-else class="space-y-4">
      <StockCard 
        v-for="stock in stocks" 
        :key="`${stock.ticker}-${stock.time}`"
        :stock="stock"
        @click="$emit('stockClick', stock)"
      />
    </div>

    <!-- Load more button -->
    <div v-if="hasMorePages && !isLoading" class="text-center mt-6">
      <button 
        @click="loadMore"
        class="btn-primary"
        :disabled="loadingMore"
      >
        <span v-if="loadingMore">Loading...</span>
        <span v-else>Load More</span>
      </button>
    </div>

    <!-- End of results -->
    <div v-if="!hasMorePages && stocks.length > 0" class="text-center mt-6">
      <p class="text-gray-500 text-sm">End of results</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useStocksStore } from '../../stores/stocks'
import StockCard from './StockCard.vue'
import ErrorMessage from '../common/ErrorMessage.vue'
import type { Stock } from '../../types/api'

// Props
interface Props {
  filters?: {
    search?: string
    rating?: string
    action?: string
    priceRange?: string
  }
}

const props = withDefaults(defineProps<Props>(), {
  filters: () => ({})
})

// Emits
const emit = defineEmits<{
  stockClick: [stock: Stock]
}>()

// Store
const stocksStore = useStocksStore()

// Local state
const loadingMore = ref(false)

// Computed
const stocks = computed(() => stocksStore.filteredStocks)
const isLoading = computed(() => stocksStore.isLoading)
const error = computed(() => stocksStore.currentError)
const hasMorePages = computed(() => stocksStore.hasMorePages)

// Methods
const loadStocks = async () => {
  await stocksStore.loadStocks()
}

const loadMore = async () => {
  if (loadingMore.value) return
  
  loadingMore.value = true
  try {
    const currentOffset = stocksStore.pagination.offset
    const limit = stocksStore.pagination.limit
    await stocksStore.loadStocks({ 
      offset: currentOffset + limit 
    })
  } catch (err) {
    console.error('Error loading more stocks:', err)
  } finally {
    loadingMore.value = false
  }
}

// Watch filters and reload
watch(() => props.filters, async (newFilters) => {
  if (newFilters.search) {
    await stocksStore.searchStocks(newFilters)
  } else {
    await loadStocks()
  }
}, { deep: true })

// Initial load
onMounted(() => {
  loadStocks()
})
</script> 