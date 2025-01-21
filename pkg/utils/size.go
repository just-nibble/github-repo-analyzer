// pkg/utils/size.go
package utils

import (
	"fmt"
	"math"
)

// BytesToHumanReadable converts bytes to human readable format
func BytesToHumanReadable(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	size := float64(bytes) / float64(div)
	return fmt.Sprintf("%.2f %cB", size, "KMGTPE"[exp])
}

// BytesToMB converts bytes to megabytes
func BytesToMB(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}

// RoundToTwoDecimals rounds a float64 to two decimal places
func RoundToTwoDecimals(num float64) float64 {
	return math.Round(num*100) / 100
}
