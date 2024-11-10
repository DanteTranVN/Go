package purchase

import "time"

func calculateSumByPeriod(purchases []Purchase, period string) int {
	now := time.Now()
	total := 0

	for _, purchase := range purchases {
		switch period {
		case "month":
			if purchase.CreatedTime.After(now.AddDate(0, -1, 0)) {
				total += purchase.Amount
			}
		case "week":
			if purchase.CreatedTime.After(now.AddDate(0, 0, -7)) {
				total += purchase.Amount
			}
		case "year":
			if purchase.CreatedTime.After(now.AddDate(-1, 0, 0)) {
				total += purchase.Amount
			}
		}
	}

	return total
}

func calculateTargetPercentage(purchases []Purchase, period string) map[string]float64 {
	total, targetSum := 0, map[string]int{}
	now := time.Now()

	for _, purchase := range purchases {
		if (period == "month" && purchase.CreatedTime.After(now.AddDate(0, -1, 0))) ||
			(period == "week" && purchase.CreatedTime.After(now.AddDate(0, 0, -7))) ||
			(period == "year" && purchase.CreatedTime.After(now.AddDate(-1, 0, 0))) {
			targetSum[purchase.Target] += purchase.Amount
			total += purchase.Amount
		}
	}

	targetPercentage := map[string]float64{}
	for target, amount := range targetSum {
		targetPercentage[target] = (float64(amount) / float64(total)) * 100
	}
	return targetPercentage
}

func calculateSumByTarget(purchases []Purchase, period string) map[string]int {
	targetTotals := map[string]int{}
	now := time.Now()

	for _, purchase := range purchases {
		if (period == "month" && purchase.CreatedTime.After(now.AddDate(0, -1, 0))) ||
			(period == "week" && purchase.CreatedTime.After(now.AddDate(0, 0, -7))) ||
			(period == "year" && purchase.CreatedTime.After(now.AddDate(-1, 0, 0))) {
			targetTotals[purchase.Target] += purchase.Amount
		}
	}
	return targetTotals
}
