// API Configuration
export const API_CONFIG = {
  BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  PUBLIC_URL: import.meta.env.VITE_API_PUBLIC_URL || 'http://localhost:8080/api/v1/public',
  ADMIN_URL: import.meta.env.VITE_API_ADMIN_URL || 'http://localhost:8080/api/v1/admin',
} as const

// API Endpoints
export const ENDPOINTS = {
  HEALTH: '/health',
  STOCKS: '/stocks',
  STOCK_DETAIL: '/stocks/:ticker',
  STOCKS_SEARCH: '/stocks/search',
  RECOMMENDATIONS: '/recommendations',
} as const

// Filter Options
export const FILTER_OPTIONS = {
  RATINGS: [
    { value: '', label: 'All Ratings' },
    { value: 'Buy', label: 'Buy' },
    { value: 'Sell', label: 'Sell' },
    { value: 'Hold', label: 'Hold' },
    { value: 'Neutral', label: 'Neutral' },
  ],
  ACTIONS: [
    { value: '', label: 'All Actions' },
    { value: 'initiated by', label: 'Initiated' },
    { value: 'upgraded by', label: 'Upgraded' },
    { value: 'downgraded by', label: 'Downgraded' },
    { value: 'target raised by', label: 'Target Raised' },
    { value: 'target lowered by', label: 'Target Lowered' },
  ],
  PRICE_RANGES: [
    { value: '', label: 'All Prices' },
    { value: '0-10', label: '$0 - $10' },
    { value: '10-50', label: '$10 - $50' },
    { value: '50-100', label: '$50 - $100' },
    { value: '100+', label: '$100+' },
  ],
} as const 