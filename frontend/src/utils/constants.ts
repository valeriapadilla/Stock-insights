export const API_CONFIG = {
  BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  PUBLIC_URL: import.meta.env.VITE_API_PUBLIC_URL || 'http://localhost:8080/api/v1/public',
  ADMIN_URL: import.meta.env.VITE_API_ADMIN_URL || 'http://localhost:8080/api/v1/admin',
} as const

export const ENDPOINTS = {
  HEALTH: '/health',
  STOCKS: '/stocks',
  STOCK_DETAIL: '/stocks/:ticker',
  STOCKS_SEARCH: '/stocks/search',
  RECOMMENDATIONS: '/recommendations',
} as const

export const FILTER_OPTIONS = {
  RATINGS: [
    { value: '', label: 'All Ratings' },
    { value: 'buy', label: 'Buy' },
    { value: 'sell', label: 'Sell' },
    { value: 'hold', label: 'Hold' },
    { value: 'neutral', label: 'Neutral' },
    { value: 'overweight', label: 'Overweight' },
    { value: 'underweight', label: 'Underweight' },
    { value: 'equal weight', label: 'Equal Weight' },
    { value: 'in-line', label: 'In-Line' }
  ]
}

export const SORT_OPTIONS = [
  { value: 'rating_to', label: 'Rating' },
  { value: 'ticker_asc', label: 'Ticker A → Z' },
  { value: 'ticker_desc', label: 'Ticker Z → A' },
  { value: 'change_percent_desc', label: 'Change % ↓ High to Low' },
  { value: 'change_percent_asc', label: 'Change % ↑ Low to High' }
] 