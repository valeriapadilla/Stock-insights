export interface SortParams {
  sort_by: string
  order: 'asc' | 'desc'
}

export function parseSortValue(sortValue: string): SortParams {
  if (sortValue.includes('_')) {
    const parts = sortValue.split('_')
    const order = parts[parts.length - 1] as 'asc' | 'desc'
    
    const field = parts.slice(0, -1).join('_')
    
    return {
      sort_by: field,
      order: order
    }
  }
  
  return {
    sort_by: sortValue,
    order: 'desc' // Default order
  }
}

export function isValidSortValue(sortValue: string): boolean {
  const validValues = [
    'time',
    'ticker_asc',
    'ticker_desc', 
    'rating_to',
    'change_percent_asc',
    'change_percent_desc'
  ]
  
  return validValues.includes(sortValue)
}

export function getDefaultSortParams(): SortParams {
  return {
    sort_by: 'time',
    order: 'desc'
  }
} 