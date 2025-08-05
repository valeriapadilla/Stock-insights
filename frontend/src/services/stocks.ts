import apiClient from './api'
import type { 
  StocksResponse, 
  StockDetailResponse, 
  StocksSearchResponse,
  Stock 
} from '../types/api'

export interface StocksParams {
  limit?: number
  offset?: number
  sort?: string
  order?: 'asc' | 'desc'
}

export interface StocksSearchParams {
  ticket?: string
  date_from?: string
  date_to?: string
  min_price?: number
  max_price?: number
  limit?: number
  offset?: number
}

export class StocksService {
  static async getStocks(params: StocksParams = {}): Promise<StocksResponse> {
    const { limit = 50, offset = 0, sort = 'time', order = 'desc' } = params
    
    const response = await apiClient.get('/stocks', {
      params: { limit, offset, sort, order }
    })
    
    return response.data
  }

  static async getStock(ticker: string): Promise<Stock> {
    const response = await apiClient.get<StockDetailResponse>(`/stocks/${ticker}`)
    return response.data.stock
  }

  static async searchStocks(params: StocksSearchParams = {}): Promise<StocksSearchResponse> {
    const { 
      ticket, 
      date_from, 
      date_to, 
      min_price, 
      max_price, 
      limit = 50, 
      offset = 0 
    } = params
    
    const response = await apiClient.get('/stocks/search', {
      params: { 
        ticket, 
        date_from, 
        date_to, 
        min_price, 
        max_price, 
        limit, 
        offset 
      }
    })
    
    return response.data
  }

  static async getStocksCount(): Promise<number> {
    const response = await apiClient.get('/stocks', {
      params: { limit: 1, offset: 0 }
    })
    
    return response.data.pagination.total
  }

  static async getStocksByAction(action: string): Promise<number> {
    const response = await apiClient.get('/stocks', {
      params: { 
        action,
        limit: 1, 
        offset: 0 
      }
    })
    
    return response.data.pagination.total
  }
} 