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
    action: string
    priceRange: string
    search: string
    sortBy: string
    sortOrder: 'asc' | 'desc'
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
      action: '',
      priceRange: '',
      search: '',
      sortBy: 'time',
      sortOrder: 'desc'
    }
  }),

  getters: {
    filteredStocks: (state) => {
      let filtered = [...state.stocks]
      
      // Filter by search
      if (state.filters.search) {
        const searchLower = state.filters.search.toLowerCase()
        filtered = filtered.filter(stock => 
          stock.ticker.toLowerCase().includes(searchLower) ||
          stock.company.toLowerCase().includes(searchLower)
        )
      }
      
      // Filter by rating
      if (state.filters.rating) {
        filtered = filtered.filter(stock => 
          stock.rating_to.toLowerCase() === state.filters.rating.toLowerCase()
        )
      }
      
      // Filter by action
      if (state.filters.action) {
        filtered = filtered.filter(stock => 
          stock.action.toLowerCase().includes(state.filters.action.toLowerCase())
        )
      }
      
      return filtered
    },
    
    hasMorePages: (state) => state.pagination.has_next,
    
    totalStocks: (state) => state.pagination.total,
    
    isLoading: (state) => state.loading,
    
    currentError: (state) => state.error
  },

  actions: {
    async loadStocks(params: { limit?: number; offset?: number; sort?: string; order?: 'asc' | 'desc' } = {}) {
      this.loading = true
      this.error = null
      
      try {
        const response: StocksResponse = await StocksService.getStocks({
          limit: params.limit || this.pagination.limit,
          offset: params.offset || this.pagination.offset,
          sort: params.sort || this.filters.sortBy,
          order: params.order || this.filters.sortOrder
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
      date_from?: string
      date_to?: string
      min_price?: number
      max_price?: number
    } = {}) {
      this.loading = true
      this.error = null
      
      try {
        const response: StocksSearchResponse = await StocksService.searchStocks({
          ...searchParams,
          limit: this.pagination.limit,
          offset: this.pagination.offset
        })
        
        this.stocks = response.stocks
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
        action: '',
        priceRange: '',
        search: '',
        sortBy: 'time',
        sortOrder: 'desc'
      }
    }
  }
}) 