package utils

func BytesToGiB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}
