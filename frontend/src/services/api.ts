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
    // Aquí podemos agregar headers de autenticación si es necesario
    console.log(`Making request to: ${config.url}`)
    return config
  },
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    console.log(`Response from: ${response.config.url}`, response.status)
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