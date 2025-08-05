import apiClient from './api'
import type { RecommendationsResponse, Recommendation } from '../types/api'

export interface RecommendationsParams {
  limit?: number
}

export class RecommendationsService {
  static async getRecommendations(params: RecommendationsParams = {}): Promise<RecommendationsResponse> {
    const { limit = 30 } = params
    
    const response = await apiClient.get('/recommendations', {
      params: { limit }
    })
    
    return response.data
  }

  static async getRecommendation(id: string): Promise<Recommendation> {
    const response = await apiClient.get(`/recommendations/${id}`)
    return response.data
  }

  
  static async getRecommendationsCount(): Promise<number> {
    try {
      const response = await apiClient.get('/recommendations', {
        params: { limit: 30 }
      })
      
      return response.data.total || 0
    } catch (error) {
      console.error('Error getting recommendations count:', error)
      return 0
    }
  }
} 