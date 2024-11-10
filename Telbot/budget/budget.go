package budget

import (
	"encoding/csv"
	"fmt"
	"github.com/tucnak/telebot"
	"os"
	"strconv"
	"strings"
)

type Budget struct {
	IDTele    int     // The user’s Telegram ID
	Category  string  // The category for the budget (e.g., "coffee")
	Amount    int     // The budgeted amount
	Duration  string  // The period for the budget (e.g., "week", "month")
	Threshold float64 // The alert threshold (e.g., 0.7 for 70%)
}

// SaveBudget saves a new budget to the CSV file without using a global variable.
func SaveBudget(budget Budget) error {
	file, err := os.OpenFile("budgets.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		strconv.Itoa(budget.IDTele),
		budget.Category,
		strconv.Itoa(budget.Amount),
		budget.Duration,
		fmt.Sprintf("%.2f", budget.Threshold),
	}

	return writer.Write(record)
}

// LoadBudgets loads all budgets from the CSV file for a specific user.
func LoadBudgets(userID int) ([]Budget, error) {
	file, err := os.Open("budgets.csv")
	if err != nil {
		if os.IsNotExist(err) {
			return []Budget{}, nil // No budgets recorded yet
		}
		return nil, err
	}
	defer file.Close()

	var budgets []Budget
	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		idTele, _ := strconv.Atoi(record[0])
		if idTele != userID {
			continue // Skip budgets not matching the specified user ID
		}

		amount, _ := strconv.Atoi(record[2])
		threshold, _ := strconv.ParseFloat(record[4], 64)

		budget := Budget{
			IDTele:    idTele,
			Category:  record[1],
			Amount:    amount,
			Duration:  record[3],
			Threshold: threshold,
		}

		budgets = append(budgets, budget)
	}

	return budgets, nil
}

func SetBudget(bot *telebot.Bot) {
	bot.Handle("/setbudget", func(m *telebot.Message) {
		args := strings.SplitN(m.Payload, " ", 3)
		if len(args) < 3 {
			bot.Send(m.Chat, "Usage: /setbudget [amount with K or M] [category] [duration (e.g., week)]")
			return
		}

		// Parse the amount with suffix
		amountString := strings.ToUpper(args[0])
		multiplier := 1

		if strings.HasSuffix(amountString, "K") {
			multiplier = 1000
			amountString = strings.TrimSuffix(amountString, "K")
		} else if strings.HasSuffix(amountString, "M") {
			multiplier = 1000000
			amountString = strings.TrimSuffix(amountString, "M")
		}

		amount, err := strconv.Atoi(amountString)
		if err != nil {
			bot.Send(m.Chat, "Please provide a valid number for the amount.")
			return
		}
		amount *= multiplier

		// Parse the category and duration
		category := args[1]
		duration := args[2]

		// Create a Budget entry for the user
		newBudget := budget.Budget{
			IDTele:    m.Sender.ID,
			Category:  category,
			Amount:    amount,
			Duration:  duration,
			Threshold: 0.7, // Set to 70%
		}

		// Save the budget
		if err := budget.SaveBudget(newBudget); err != nil {
			bot.Send(m.Chat, "Failed to save budget.")
			return
		}

		bot.Send(m.Chat, fmt.Sprintf("Budget set for %s: %s every %s.", category, formatNumber(amount), duration))
	})
}
