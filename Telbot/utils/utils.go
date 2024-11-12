package utils

import (
	"fmt"
)

func FormatNumber(amount int) string {
	millions := amount / 1_000_000
	remainderAfterMillions := amount % 1_000_000
	thousands := remainderAfterMillions / 1_000

	switch {
	case millions > 0 && thousands > 0:
		return fmt.Sprintf("%dM and %dK", millions, thousands)
	case millions > 0:
		return fmt.Sprintf("%dM", millions)
	case thousands > 0:
		return fmt.Sprintf("%dK", thousands)
	default:
		return fmt.Sprintf("%d", amount)
	}
}
