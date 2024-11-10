package purchase

import (
	"fmt"
	"github.com/tucnak/telebot"
	"os"
	"strconv"
	"strings"
	"time"
)

func formatNumber(amount int) string {
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

// Purchase struct to store purchase information with CreatedTime
type Purchase struct {
	IDTele      int
	AccountName string
	Amount      int
	Target      string
	CreatedTime time.Time
}

// RegisterHandlers registers purchase-related handlers to the bot
func RegisterHandlers(bot *telebot.Bot) {
	// Handle the /purchase command
	bot.Handle("/purchase", func(m *telebot.Message) {
		args := strings.SplitN(m.Payload, " ", 2)
		if len(args) < 2 {
			bot.Send(m.Chat, "Usage: /purchase [amount with K or M] [target]")
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

		// Parse the target (e.g., "education")
		target := args[1]

		// Create a Purchase entry with the current timestamp
		purchase := Purchase{
			IDTele:      m.Sender.ID,
			AccountName: m.Sender.Username,
			Amount:      amount,
			Target:      target,
			CreatedTime: time.Now(), // Capture the current time
		}

		// Save the purchase record to a file
		if err := savePurchaseToFile(purchase); err != nil {
			bot.Send(m.Chat, "Failed to save purchase record.")
			return
		}

		bot.Send(m.Chat, fmt.Sprintf("Recorded purchase: %d for %s.", purchase.Amount, purchase.Target))
	})
}

func savePurchaseToFile(purchase Purchase) error {
	// Open or create the file in append mode
	file, err := os.OpenFile("purchase_records.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write purchase data as a CSV line, using YYYY-MM-DD format for CreatedTime
	_, err = file.WriteString(fmt.Sprintf("%d,%s,%d,%s,%s\n", purchase.IDTele, purchase.AccountName, purchase.Amount, purchase.Target, purchase.CreatedTime.Format("2006-01-02")))
	return err
}

func RegisterReportCommands(bot *telebot.Bot) {
	// Command to get total sum by period
	bot.Handle("/sumPurchases", func(m *telebot.Message) {
		purchases, err := loadPurchases()
		if err != nil {
			bot.Send(m.Chat, "Failed to load purchase records.")
			return
		}

		lastMonth := calculateSumByPeriod(purchases, "month")
		lastWeek := calculateSumByPeriod(purchases, "week")
		lastYear := calculateSumByPeriod(purchases, "year")

		message := fmt.Sprintf("Sum of Purchases:\nLast Month: %d\nLast Week: %d\nLast Year: %d", lastMonth, lastWeek, lastYear)
		bot.Send(m.Chat, message)
	})

	// Command to get target percentage by period
	bot.Handle("/targetPercentage", func(m *telebot.Message) {
		purchases, err := loadPurchases()
		if err != nil {
			bot.Send(m.Chat, "Failed to load purchase records.")
			return
		}

		period := "month" // This could be modified to take arguments for different periods
		targetData := calculateTargetPercentage(purchases, period)

		message := "Target Percentage in " + period + ":\n"
		for target, percentage := range targetData {
			message += fmt.Sprintf("%s: %.2f%%\n", target, percentage)
		}
		bot.Send(m.Chat, message)
	})

	// Command to get sum by target
	bot.Handle("/sumByTarget", func(m *telebot.Message) {
		purchases, err := loadPurchases()
		if err != nil {
			bot.Send(m.Chat, "Failed to load purchase records.")
			return
		}

		period := "month" // This could be modified to take arguments for different periods
		targetTotals := calculateSumByTarget(purchases, period)

		message := fmt.Sprintf("Sum by Target in %s:\n", period)
		for target, sum := range targetTotals {
			message += fmt.Sprintf("%s: %d\n", target, sum)
		}
		bot.Send(m.Chat, message)
	})
	bot.Handle("/targetSummary", func(m *telebot.Message) {
		purchases, err := loadPurchases()
		if err != nil {
			bot.Send(m.Chat, "Failed to load purchase records.")
			return
		}

		// Define the period (e.g., "month")
		period := "month" // Adjust this or modify to take arguments for different periods

		// Calculate sum by target
		targetTotals := calculateSumByTarget(purchases, period)

		// Calculate target percentage
		targetPercentages := calculateTargetPercentage(purchases, period)

		// Create a response message combining both totals and percentages
		message := fmt.Sprintf("Target Summary in %s:\n", period)
		for target, sum := range targetTotals {
			percentage := targetPercentages[target]
			message += fmt.Sprintf("%s: %s (%.2f%%)\n", target, formatNumber(sum), percentage)
		}

		bot.Send(m.Chat, message)
	})

}
