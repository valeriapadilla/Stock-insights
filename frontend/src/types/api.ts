export interface ApiResponse<T> {
  data?: T
  error?: string
  message?: string
}

export interface PaginationResponse {
  total: number
  limit: number
  offset: number
  has_next: boolean
}

export interface Stock {
  ticker: string
  company: string
  target_from: string
  target_to: string
  rating_from: string
  rating_to: string
  action: string
  brokerage: string
  time: string 
  created_at: string 
  updated_at: string
  change_percent?: string
}

// /api/v1/public/stocks
export interface StocksResponse {
  stocks: Stock[]
  pagination: PaginationResponse
}

// /api/v1/public/stocks/:ticker
export interface StockDetailResponse {
  stock: Stock
}

// /api/v1/public/stocks/search
export interface StocksSearchResponse {
  stocks: Stock[]
  pagination: PaginationResponse
  filters_applied: {
    ticket?: string
    rating?: string
    sort_by?: string
    order?: string
  }
}

export interface Recommendation {
  id: string
  ticker: string
  score: number 
  explanation: string
  run_at: string 
  rank: number
}

// /api/v1/public/recommendations
export interface RecommendationsResponse {
  recommendations: Recommendation[]
  total: number
  limit: number
}

export interface FilterOption {
  value: string
  label: string
}

export interface FilterState {
  rating: string
  search: string
  sort_by: string
  order: 'asc' | 'desc'
}

export interface DashboardStats {
  totalStocks: number
  upgradesToday: number
  downgradesToday: number
  topRecommendations: number
}

export interface ApiError {
  error: string
  message: string
} 