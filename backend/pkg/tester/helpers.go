package tester

// Helper to calculate the max of a slice
func max(slice []float64) float64 {
	if len(slice) == 0 {
		return 0
	}
	max := slice[0]
	for _, v := range slice {
		if v > max {
			max = v
		}
	}
	return max
}

// Helper to calculate percentile
func percentile(sortedSlice []float64, percent float64) float64 {
	if len(sortedSlice) == 0 {
		return 0
	}
	rank := (percent / 100) * float64(len(sortedSlice)-1)
	lower := int(rank)
	upper := lower + 1
	if upper >= len(sortedSlice) {
		return sortedSlice[lower]
	}
	fraction := rank - float64(lower)
	return sortedSlice[lower] + fraction*(sortedSlice[upper]-sortedSlice[lower])
}
