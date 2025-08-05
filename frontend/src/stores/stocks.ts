import { defineStore } from 'pinia'
import { StocksService } from '../services/stocks'
import type { Stock, StocksResponse, StocksSearchResponse } from '../types/api'

interface StocksState {
  stocks: Stock[]
  currentStock: Stock | null
  loading: boolean
  error: string | null
  pagination: {
    total: number
    limit: number
    offset: number
    has_next: boolean
  }
  filters: {
    rating: string
    search: string
    sort_by: string
    order: 'asc' | 'desc'
  }
}

export const useStocksStore = defineStore('stocks', {
  state: (): StocksState => ({
    stocks: [],
    currentStock: null,
    loading: false,
    error: null,
    pagination: {
      total: 0,
      limit: 50,
      offset: 0,
      has_next: false
    },
    filters: {
      rating: '',
      search: '',
      sort_by: 'time',
      order: 'desc'
    }
  }),

  getters: {
    // Remove filteredStocks getter - always use backend filtering
    hasMorePages: (state) => state.pagination.has_next,
    
    totalStocks: (state) => state.pagination.total,
    
    isLoading: (state) => state.loading,
    
    currentError: (state) => state.error
  },

  actions: {
    async loadStocks(params: { limit?: number; offset?: number; sort_by?: string; order?: 'asc' | 'desc' } = {}) {
      this.loading = true
      this.error = null
      
      try {
        const response: StocksResponse = await StocksService.getStocks({
          limit: params.limit || this.pagination.limit,
          offset: params.offset || this.pagination.offset,
          sort_by: params.sort_by || this.filters.sort_by,
          order: params.order || this.filters.order
        })
        
        // If offset is 0, replace stocks, otherwise append
        if (params.offset === 0 || !params.offset) {
          this.stocks = response.stocks
        } else {
          this.stocks = [...this.stocks, ...response.stocks]
        }
        
        this.pagination = response.pagination
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Error loading stocks'
        console.error('Error loading stocks:', error)
      } finally {
        this.loading = false
      }
    },

    async searchStocks(searchParams: {
      ticket?: string
      rating?: string
      sort_by?: string
      order?: 'asc' | 'desc'
      limit?: number
      offset?: number
    } = {}) {
      this.loading = true
      this.error = null
      
      try {
        const response: StocksSearchResponse = await StocksService.searchStocks({
          ...searchParams,
          limit: searchParams.limit || this.pagination.limit,
          offset: searchParams.offset || this.pagination.offset
        })
        
        // If offset is 0, replace stocks, otherwise append
        if (searchParams.offset === 0 || !searchParams.offset) {
          this.stocks = response.stocks
        } else {
          this.stocks = [...this.stocks, ...response.stocks]
        }
        
        this.pagination = response.pagination
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Error searching stocks'
        console.error('Error searching stocks:', error)
      } finally {
        this.loading = false
      }
    },

    async loadStock(ticker: string) {
      this.loading = true
      this.error = null
      
      try {
        this.currentStock = await StocksService.getStock(ticker)
      } catch (error) {
        this.error = error instanceof Error ? error.message : 'Error loading stock'
        console.error('Error loading stock:', error)
      } finally {
        this.loading = false
      }
    },

    updateFilters(newFilters: Partial<StocksState['filters']>) {
      this.filters = { ...this.filters, ...newFilters }
    },

    clearError() {
      this.error = null
    },

    reset() {
      this.stocks = []
      this.currentStock = null
      this.loading = false
      this.error = null
      this.pagination = {
        total: 0,
        limit: 50,
        offset: 0,
        has_next: false
      }
      this.filters = {
        rating: '',
        search: '',
        sort_by: 'time',
        order: 'desc'
      }
    }
  }
}) 