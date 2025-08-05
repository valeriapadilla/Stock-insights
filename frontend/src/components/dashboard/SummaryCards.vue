<template>
  <div class="summary-cards">

    <div v-if="isLoading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div v-for="i in 4" :key="i" class="bg-gray-800 border border-gray-700 rounded-lg p-6 animate-pulse">
        <div class="h-4 bg-gray-700 rounded mb-2"></div>
        <div class="h-8 bg-gray-700 rounded mb-2"></div>
        <div class="h-3 bg-gray-700 rounded"></div>
      </div>
    </div>

    <ErrorMessage 
      v-else-if="error"
      :message="error"
      :showRetry="true"
      @retry="loadDashboardData"
    />

    <!-- Data state -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <SummaryCard
        title="Total Stocks"
        :value="dashboardStats.totalStocks"
        description="Active stocks tracked"
        type="default"
      />

      <SummaryCard
        title="Upgrades Today"
        :value="dashboardStats.upgradesToday"
        description="+12% from yesterday"
        type="success"
        :showChange="true"
        changeType="up"
        changeText="+12% from yesterday"
      />

      <SummaryCard
        title="Downgrades Today"
        :value="dashboardStats.downgradesToday"
        description="-8% from yesterday"
        type="error"
        :showChange="true"
        changeType="down"
        changeText="-8% from yesterday"
      />

      <SummaryCard
        title="Top Recommendations"
        :value="dashboardStats.topRecommendations"
        description="Updated daily"
        type="default"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useStocksStore } from '../../stores/stocks'
import { useRecommendationsStore } from '../../stores/recommendations'
import SummaryCard from './SummaryCard.vue'
import ErrorMessage from '../common/ErrorMessage.vue'
import LoadingSpinner from '../common/LoadingSpinner.vue'

// Stores
const stocksStore = useStocksStore()
const recommendationsStore = useRecommendationsStore()

const isLoading = ref(false)
const error = ref('')

const dashboardStats = computed(() => ({
  totalStocks: stocksStore.totalStocks,
  upgradesToday: 127, // TODO: Get from API
  downgradesToday: 43, // TODO: Get from API
  topRecommendations: recommendationsStore.totalRecommendations
}))

const loadDashboardData = async () => {
  isLoading.value = true
  error.value = ''

  try {
    await Promise.all([
      stocksStore.loadStocks({ limit: 1 }), // Get total count
      recommendationsStore.loadRecommendations({ limit: 30 }) // Get total count
    ])
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load dashboard data'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadDashboardData()
})
</script> 