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