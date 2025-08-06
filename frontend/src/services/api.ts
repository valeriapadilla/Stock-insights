import axios from 'axios'
import type { AxiosInstance, AxiosResponse } from 'axios'
import { API_CONFIG } from '../utils/constants'
import type { ApiError } from '../types/api'

const apiClient: AxiosInstance = axios.create({
  baseURL: API_CONFIG.PUBLIC_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

apiClient.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    return response
  },
  (error) => {
    console.error('Response error:', error)
    
    if (error.response) {
      const apiError: ApiError = error.response.data
      console.error('API Error:', apiError)
    }
    
    return Promise.reject(error)
  }
)

export default apiClient 