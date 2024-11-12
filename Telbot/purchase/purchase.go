package purchase

import (
	"Telbot/utils"
	"encoding/csv"
	"fmt"
	"github.com/tucnak/telebot"
	"os"
	"strconv"
	"strings"
	"time"
)

// Purchase struct to store purchase information with CreatedTime
type Purchase struct {
	IDTele      int
	AccountName string
	Amount      int
	Target      string
	CreatedTime time.Time
}

type Budget struct {
	IDTele    int     // The user’s Telegram ID
	Category  string  // The category for the budget (e.g., "coffee")
	Amount    int     // The budgeted amount
	Duration  string  // The period for the budget (e.g., "week", "month")
	Threshold float64 // The alert threshold (e.g., 0.7 for 70%)
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

func CalculateSpent(userID int, target string, period string) int {
	purchases, err := loadPurchases()
	if err != nil {
		// Trả về 0 nếu không thể tải dữ liệu mua hàng
		return 0
	}

	totalSpent := 0
	now := time.Now()

	for _, purchase := range purchases {
		// Lọc theo ID người dùng và mục tiêu
		if purchase.IDTele == userID && purchase.Target == target {
			// Kiểm tra xem giao dịch có thuộc khoảng thời gian chỉ định không
			switch period {
			case "week":
				// Kiểm tra nếu giao dịch thuộc tuần hiện tại
				weekAgo := now.AddDate(0, 0, -7)
				if purchase.CreatedTime.After(weekAgo) {
					totalSpent += purchase.Amount
				}
			case "month":
				// Kiểm tra nếu giao dịch thuộc tháng hiện tại
				monthAgo := now.AddDate(0, -1, 0)
				if purchase.CreatedTime.After(monthAgo) {
					totalSpent += purchase.Amount
				}
			case "year":
				// Kiểm tra nếu giao dịch thuộc năm hiện tại
				yearAgo := now.AddDate(-1, 0, 0)
				if purchase.CreatedTime.After(yearAgo) {
					totalSpent += purchase.Amount
				}
			}
		}
	}

	return totalSpent
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
		message, _ := CheckBudgetAlert(m.Sender.ID)
		if message != "" {

			bot.Send(m.Chat, message)
		}
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
		// Lấy tham số period từ tin nhắn người dùng
		period := "month" // Mặc định là "month"
		args := strings.Split(m.Text, " ")
		if len(args) > 1 {
			period = args[1]
			if period != "month" && period != "year" && period != "week" {
				bot.Send(m.Chat, "Please specify a valid period: 'month' or 'year'.")
				return
			}
		}

		// Tải dữ liệu mua hàng
		purchases, err := loadPurchases()
		if err != nil {
			bot.Send(m.Chat, "Failed to load purchase records.")
			return
		}

		// Tính tổng theo mục tiêu (target)
		targetTotals := calculateSumByTarget(purchases, period)

		// Tính phần trăm theo mục tiêu (target)
		targetPercentages := calculateTargetPercentage(purchases, period)

		// Tạo tin nhắn phản hồi, bao gồm cả tổng và phần trăm cho từng mục tiêu
		message := fmt.Sprintf("Target Summary in %s:\n", period)
		for target, sum := range targetTotals {
			percentage := targetPercentages[target]
			message += fmt.Sprintf("%s: %s (%.2f%%)\n", target, utils.FormatNumber(sum), percentage)
		}

		// Gửi tin nhắn phản hồi
		bot.Send(m.Chat, message)
	})
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

func SetBudget(bot *telebot.Bot) {
	bot.Handle("/setBudget", func(m *telebot.Message) {
		args := strings.SplitN(m.Payload, " ", 3)
		if len(args) < 3 {
			bot.Send(m.Chat, "Usage: /setBudget [amount with K or M] [category] [duration (e.g., week)]")
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
		newBudget := Budget{
			IDTele:    m.Sender.ID,
			Category:  category,
			Amount:    amount,
			Duration:  duration,
			Threshold: 0.7, // Set to 70%
		}

		// Save the budget
		if err := SaveBudget(newBudget); err != nil {
			bot.Send(m.Chat, "Failed to save budget.")
			return
		}

		bot.Send(m.Chat, fmt.Sprintf("Budget set for %s: %s every %s.", category, utils.FormatNumber(amount), duration))
	})
}

// Hàm xử lý xem ngân sách của người dùng
func ViewBudget(bot *telebot.Bot) {
	bot.Handle("/viewBudget", func(m *telebot.Message) {
		budgets, err := LoadBudgets(m.Sender.ID)
		if err != nil {
			bot.Send(m.Chat, "Failed to load budgets.")
			return
		}

		if len(budgets) == 0 {
			bot.Send(m.Chat, "You have no budgets set.")
			return
		}

		message := "Your Budgets:\n"
		for _, budget := range budgets {
			message += fmt.Sprintf("- %s: %s every %s\n", budget.Category, utils.FormatNumber(budget.Amount), budget.Duration)
		}

		bot.Send(m.Chat, message)
	})
}

// Hàm kiểm tra mức độ sử dụng ngân sách và cảnh báo khi gần vượt quá ngưỡng
func CheckBudget(bot *telebot.Bot) {
	bot.Handle("/checkBudget", func(m *telebot.Message) {
		budgets, err := LoadBudgets(m.Sender.ID)
		if err != nil {
			bot.Send(m.Chat, "Failed to load budgets.")
			return
		}

		if len(budgets) == 0 {
			bot.Send(m.Chat, "You have no budgets set.")
			return
		}

		message := "Budget Usage:\n"
		for _, budget := range budgets {
			spent := CalculateSpent(m.Sender.ID, budget.Category, budget.Duration) // Hàm để tính chi tiêu theo Category và Duration
			percentage := float64(spent) / float64(budget.Amount)

			message += fmt.Sprintf("%s: Spent %s of %s (%.2f%%)\n", budget.Category, utils.FormatNumber(spent), utils.FormatNumber(budget.Amount), percentage*100)

			if percentage >= budget.Threshold {
				message += fmt.Sprintf("⚠️ Warning: You've spent %.2f%% of your %s budget!\n", percentage*100, budget.Category)
			}
		}

		bot.Send(m.Chat, message)
	})
}

func CheckBudgetAlert(Id int) (string, error) {
	budgets, err := LoadBudgets(Id)
	if err != nil {
		return "", nil
	}

	if len(budgets) == 0 {
		return "", nil
	}

	for _, budget := range budgets {
		// Tính tổng chi tiêu cho mục tiêu và khoảng thời gian này
		spent := CalculateSpent(Id, budget.Category, budget.Duration)
		percentageSpent := float64(spent) / float64(budget.Amount)

		// Nếu chi tiêu vượt quá 50%, gửi cảnh báo
		if percentageSpent >= 0.5 {
			remaining := budget.Amount - spent
			message := fmt.Sprintf("⚠️ Alert: You've spent %.2f%% of your %s budget for %s.\nRemaining amount: %s",
				percentageSpent*100, budget.Duration, budget.Category, utils.FormatNumber(remaining))
			return message, nil
		}
	}
	return "", nil
}
