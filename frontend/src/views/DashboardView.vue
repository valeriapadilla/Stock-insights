<template>
  <div class="min-h-screen bg-gray-900">
    <Header />
    
    <div class="max-w-7xl mx-auto px-6 py-8">
      <div class="mb-8">
        <SummaryCards />
      </div>

      <NavigationTabs 
        v-model:activeTab="activeTab"
        @update:activeTab="handleTabChange"
      />

      <div class="content-area">
        <div v-if="activeTab === 'stocks'" class="space-y-6">
          <div class="bg-gray-800 border border-gray-700 rounded-lg p-6">
            <h2 class="text-xl font-semibold text-white mb-2">Stock Analysis</h2>
            <p class="text-gray-400 mb-6">Search and filter through all tracked stocks.</p>
            
            <FilterBar 
              v-model="stockFilters"
              @search="handleStockSearch"
            />
          </div>

          <div class="space-y-4">
            <div v-if="stocksStore.isLoading" class="flex justify-center py-8">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-green-500"></div>
            </div>
            
            <div v-else-if="stocksStore.currentError" class="text-center py-8">
              <p class="text-red-400">{{ stocksStore.currentError }}</p>
              <button 
                @click="loadStocks"
                class="mt-4 px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600"
              >
                Retry
              </button>
            </div>
            
            <div v-else class="grid gap-4">
              <StockCard 
                v-for="stock in stocksStore.stocks" 
                :key="stock.ticker"
                :stock="stock"
              />
            </div>

            <div v-if="stocksStore.hasMorePages" class="flex justify-center mt-8">
              <button 
                @click="loadMoreStocks"
                class="px-6 py-2 bg-gray-700 text-white rounded-md hover:bg-gray-600"
              >
                Load More
              </button>
            </div>
          </div>
        </div>

        <div v-else-if="activeTab === 'recommendations'" class="space-y-6">
          <div class="bg-gray-800 border border-gray-700 rounded-lg p-6">
            <h2 class="text-xl font-semibold text-white mb-2">Daily Recommendations</h2>
            <p class="text-gray-400 mb-6">Top stock recommendations based on our algorithm.</p>
          </div>

          <!-- Recommendations List -->
          <div class="space-y-4">
            <div v-if="recommendationsStore.isLoading" class="flex justify-center py-8">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-green-500"></div>
            </div>
            
            <div v-else-if="recommendationsStore.currentError" class="text-center py-8">
              <p class="text-red-400">{{ recommendationsStore.currentError }}</p>
              <button 
                @click="loadRecommendations"
                class="mt-4 px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600"
              >
                Retry
              </button>
            </div>
            
            <div v-else class="grid gap-4">
              <RecommendationCard 
                v-for="recommendation in recommendationsStore.recommendations" 
                :key="recommendation.id"
                :recommendation="recommendation"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useStocksStore } from '../stores/stocks'
import { useRecommendationsStore } from '../stores/recommendations'
import Header from '../components/common/Header.vue'
import SummaryCards from '../components/dashboard/SummaryCards.vue'
import NavigationTabs from '../components/common/NavigationTabs.vue'
import FilterBar from '../components/stocks/FilterBar.vue'
import StockCard from '../components/stocks/StockCard.vue'
import RecommendationCard from '../components/recommendations/RecommendationCard.vue'

const stocksStore = useStocksStore()
const recommendationsStore = useRecommendationsStore()

const activeTab = ref<'stocks' | 'recommendations'>('stocks')
const stockFilters = ref({
  search: '',
  rating: '',
  sort_by: 'ticker_asc',
  order: 'asc' as 'asc' | 'desc'
})

const handleTabChange = (tab: 'stocks' | 'recommendations') => {
  activeTab.value = tab
  if (tab === 'stocks') {
    loadStocks()
  } else {
    loadRecommendations()
  }
}

const handleStockSearch = (query: string) => {
  stockFilters.value.search = query
  loadStocks()
}

const loadStocks = async () => {
  await stocksStore.searchStocks({
    ticket: stockFilters.value.search,
    rating: stockFilters.value.rating,
    sort_by: stockFilters.value.sort_by,
    order: stockFilters.value.order,
    limit: 20,
    offset: 0
  })
}

const loadMoreStocks = async () => {
  const currentOffset = stocksStore.pagination.offset
  await stocksStore.searchStocks({
    ticket: stockFilters.value.search,
    rating: stockFilters.value.rating,
    sort_by: stockFilters.value.sort_by,
    order: stockFilters.value.order,
    limit: 20,
    offset: currentOffset + stocksStore.pagination.limit
  })
}

const loadRecommendations = async () => {
  await recommendationsStore.loadRecommendations({
    limit: 30
  })
}

watch(stockFilters, () => {
  if (activeTab.value === 'stocks') {
    loadStocks()
  }
}, { deep: true })

onMounted(() => {
  loadStocks()
})
</script> 