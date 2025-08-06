export const getArrowPath = (changePercent?: string, targetFrom?: string, targetTo?: string) => {
  let isPositive = true
  
  if (changePercent) {
    const change = parseFloat(changePercent.replace(/[+%]/g, ''))
    isPositive = change >= 0
  } else if (targetFrom && targetTo) {
    const fromPrice = parseFloat(targetFrom.replace('$', ''))
    const toPrice = parseFloat(targetTo.replace('$', ''))
    isPositive = toPrice > fromPrice
  }
  
  if (isPositive) {
    // Arrow pointing up
    return "M3.293 9.707a1 1 0 010-1.414l6-6a1 1 0 011.414 0l6 6a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L4.707 9.707a1 1 0 01-1.414 0z"
  } else {
    // Arrow pointing down
    return "M16.707 10.293a1 1 0 010 1.414l-6 6a1 1 0 01-1.414 0l-6-6a1 1 0 011.414-1.414L9 14.586V3a1 1 0 012 0v11.586l4.293-4.293a1 1 0 011.414 0z"
  }
}

export const getRatingColor = (rating: string) => {
  switch (rating.toLowerCase()) {
    case 'buy':
      return 'bg-green-600 text-white'
    case 'sell':
      return 'bg-red-600 text-white'
    case 'hold':
      return 'bg-yellow-600 text-white'
    case 'neutral':
      return 'bg-blue-600 text-white'
    default:
      return 'bg-gray-600 text-gray-300'
  }
}

export const getChangeColor = (changePercent?: string, targetFrom?: string, targetTo?: string) => {
  if (changePercent) {
    const change = parseFloat(changePercent.replace(/[+%]/g, ''))
    return change >= 0 ? 'text-green-400' : 'text-red-400'
  }
  
  if (targetFrom && targetTo) {
    const fromPrice = parseFloat(targetFrom.replace('$', ''))
    const toPrice = parseFloat(targetTo.replace('$', ''))
    return toPrice > fromPrice ? 'text-green-400' : 'text-red-400'
  }
  
  return 'text-gray-400'
}

export const getFallbackChangePercentage = (targetFrom: string, targetTo: string) => {
  const fromPrice = parseFloat(targetFrom.replace('$', ''))
  const toPrice = parseFloat(targetTo.replace('$', ''))
  const change = ((toPrice - fromPrice) / fromPrice) * 100
  return `${change > 0 ? '+' : ''}${change.toFixed(1)}%`
}

export const formatTime = (timeString: string, includeYear = false) => {
  const date = new Date(timeString)
  const options: Intl.DateTimeFormatOptions = {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }
  
  if (includeYear) {
    options.year = 'numeric'
  }
  
  return date.toLocaleDateString('en-US', options)
} 