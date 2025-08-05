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
    { value: 'in-line', label: 'In-Line' }
  ],
  
  ACTIONS: [
    { value: '', label: 'All Actions' },
    { value: 'target_raised', label: 'Target Raised' },
    { value: 'target_lowered', label: 'Target Lowered' },
    { value: 'rating_upgraded', label: 'Rating Upgraded' },
    { value: 'rating_downgraded', label: 'Rating Downgraded' },
    { value: 'initiated', label: 'Initiated' },
    { value: 'reiterated', label: 'Reiterated' }
  ],
  
  PRICE_RANGES: [
    { value: '', label: 'All Prices' },
    { value: '0-10', label: '$0 - $10' },
    { value: '10-25', label: '$10 - $25' },
    { value: '25-50', label: '$25 - $50' },
    { value: '50-100', label: '$50 - $100' },
    { value: '100-250', label: '$100 - $250' },
    { value: '250+', label: '$250+' }
  ]
}

export const SORT_OPTIONS = [
  { value: 'time', label: 'Time' },
  { value: 'ticker', label: 'Ticker' },
  { value: 'company', label: 'Company' },
  { value: 'target_to', label: 'Target Price' },
  { value: 'rating_to', label: 'Rating' }
] 