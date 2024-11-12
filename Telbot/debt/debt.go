package debt

import (
	"Telbot/utils"
	"encoding/csv"
	"fmt"
	"github.com/tucnak/telebot"
	"os"
	"strconv"
	"strings"
)

var debtorMap = make(map[string]int)

func LoadDebtRecords() error {
	file, err := os.Open("debtors.csv")
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Nếu file không tồn tại, không có lỗi
		}
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		name := record[0]
		amount, _ := strconv.Atoi(record[1])
		debtorMap[name] = amount
	}

	return nil
}

func SaveDebtRecords() error {
	file, err := os.Create("debtors.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for name, amount := range debtorMap {
		record := []string{name, strconv.Itoa(amount)}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func RegisterHandlers(bot *telebot.Bot) {
	if err := LoadDebtRecords(); err != nil {
		fmt.Println("Failed to load debt records:", err)
	}

	bot.Handle("/addDebtor", func(m *telebot.Message) {
		args := strings.Fields(m.Payload) // Chia dựa trên khoảng trắng để bỏ qua dấu cách thừa
		if len(args) < 2 {
			bot.Send(m.Chat, "Usage: /addDebtor [name] [amount with K for thousand or M for million]")
			return
		}

		name := args[0]
		amountString := strings.ToUpper(args[1])
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
		debtorMap[name] += amount

		if err := SaveDebtRecords(); err != nil {
			bot.Send(m.Chat, "Failed to save debt records.")
			return
		}

		bot.Send(m.Chat, fmt.Sprintf("Updated %s's debt to %s.", name, utils.FormatNumber(debtorMap[name])))
	})

	bot.Handle("/listDebtors", func(m *telebot.Message) {
		if len(debtorMap) == 0 {
			bot.Send(m.Chat, "No debtors recorded.")
			return
		}
		reply := "Current debt list:\n"
		for name, amount := range debtorMap {
			reply += fmt.Sprintf("%s owes %s\n", name, utils.FormatNumber(amount))
		}
		bot.Send(m.Chat, reply)
	})

	bot.Handle("/delDebtor", func(m *telebot.Message) {
		name := strings.TrimSpace(m.Payload)
		if name == "" {
			bot.Send(m.Chat, "Usage: /delDebtor [name]")
			return
		}

		if _, exists := debtorMap[name]; exists {
			delete(debtorMap, name)

			if err := SaveDebtRecords(); err != nil {
				bot.Send(m.Chat, "Failed to update debt records.")
				return
			}

			bot.Send(m.Chat, fmt.Sprintf("Deleted %s from debt records.", name))
		} else {
			bot.Send(m.Chat, fmt.Sprintf("No debtor found with the name %s.", name))
		}
	})
}
