import apiClient from './api'
import type { 
  StocksResponse, 
  StockDetailResponse, 
  StocksSearchResponse,
  Stock 
} from '../types/api'
import { parseSortValue } from '../utils/sort'

export interface StocksParams {
  limit?: number
  offset?: number
  sort_by?: string
  order?: 'asc' | 'desc'
}

export interface StocksSearchParams {
  ticket?: string
  rating?: string
  sort_by?: string
  order?: 'asc' | 'desc'
  limit?: number
  offset?: number
}

export class StocksService {
  static async getStocks(params: StocksParams = {}): Promise<StocksResponse> {
    const { limit = 50, offset = 0, sort_by = 'time', order = 'desc' } = params
    
    const response = await apiClient.get('/stocks', {
      params: { limit, offset, sort_by, order }
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
      rating,
      sort_by,
      limit = 50, 
      offset = 0 
    } = params
    
    const sortParams = sort_by ? parseSortValue(sort_by) : { sort_by: 'time', order: 'desc' }
    
    const response = await apiClient.get('/stocks/search', {
      params: { 
        ticket, 
        rating,
        sort_by: sortParams.sort_by,
        order: sortParams.order,
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
} 