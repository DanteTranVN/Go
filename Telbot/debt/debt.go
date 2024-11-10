package debt

import (
	"fmt"
	"github.com/tucnak/telebot"
	"strconv"
	"strings"
)

var debtorMap = make(map[string]int)

// RegisterHandlers registers debt-related handlers to the bot
func RegisterHandlers(bot *telebot.Bot) {
	// Add or update a debtor's debt
	bot.Handle("/addDebtor", func(m *telebot.Message) {
		args := strings.SplitN(m.Payload, " ", 2)
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
		bot.Send(m.Chat, fmt.Sprintf("Updated %s's debt to %d.", name, debtorMap[name]))
	})

	// List all debtors
	bot.Handle("/listDebtors", func(m *telebot.Message) {
		if len(debtorMap) == 0 {
			bot.Send(m.Chat, "No debtors recorded.")
			return
		}
		reply := "Current debt list:\n"
		for name, amount := range debtorMap {
			reply += fmt.Sprintf("%s owes %d\n", name, amount)
		}
		bot.Send(m.Chat, reply)
	})

	// Delete a debtor
	bot.Handle("/delDebtor", func(m *telebot.Message) {
		name := strings.TrimSpace(m.Payload)
		if name == "" {
			bot.Send(m.Chat, "Usage: /delDebtor [name]")
			return
		}

		if _, exists := debtorMap[name]; exists {
			delete(debtorMap, name)
			bot.Send(m.Chat, fmt.Sprintf("Deleted %s from debt records.", name))
		} else {
			bot.Send(m.Chat, fmt.Sprintf("No debtor found with the name %s.", name))
		}
	})
}
