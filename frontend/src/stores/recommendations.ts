import { defineStore } from 'pinia'
import { RecommendationsService } from '../services/recommendations'
import type { Recommendation, RecommendationsResponse } from '../types/api'

interface RecommendationsState {
  recommendations: Recommendation[]
  loading: boolean
  error: string | null
  total: number
  limit: number
}

export const useRecommendationsStore = defineStore('recommendations', {
  state: (): RecommendationsState => ({
    recommendations: [],
    loading: false,
    error: null,
    total: 0,
    limit: 30
  }),

  getters: {
    sortedRecommendations: (state) => {
      return [...state.recommendations].sort((a, b) => b.score - a.score)
    },
    
    topRecommendations: (state) => {
      return [...state.recommendations]
        .sort((a, b) => b.score - a.score)
        .slice(0, 5)
    },
    
    isLoading: (state) => state.loading,
    
    currentError: (state) => state.error,
    
    totalRecommendations: (state) => state.total
  },

  actions: {
    async loadRecommendations(params: { limit?: number } = {}) {
      this.loading = true
      this.error = null
      
      try {
        const response: RecommendationsResponse = await RecommendationsService.getRecommendations({ 
          limit: params.limit || this.limit 
        })
        
        this.recommendations = response.recommendations
        this.total = response.total
        this.limit = response.limit
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Error loading recommendations'
        console.error('Error loading recommendations:', error)
      } finally {
        this.loading = false
      }
    },

    async getRecommendationsCount(): Promise<number> {
      try {
        return await RecommendationsService.getRecommendationsCount()
      } catch (error) {
        console.error('Error getting recommendations count:', error)
        return 0
      }
    },

    clearError() {
      this.error = null
    },

    reset() {
      this.recommendations = []
      this.loading = false
      this.error = null
      this.total = 0
      this.limit = 10
    }
  }
}) 