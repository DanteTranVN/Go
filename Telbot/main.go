package main

import (
	"Telbot/debt"
	"Telbot/purchase"
	"github.com/tucnak/telebot"
	"log"
	"time"
)

func main() {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  "7821718420:AAFabsK75GWyGvD2yq45YdFfgu3ffvwynHM",
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	// Register handlers from each package
	debt.RegisterHandlers(bot)
	purchase.RegisterHandlers(bot)       // Handles purchase-related commands
	purchase.RegisterReportCommands(bot) // Handles reporting-related commands
	//saving.RegisterHandlers(bot)
	purchase.SetBudget(bot)
	purchase.CheckBudget(bot)
	purchase.ViewBudget(bot)
	bot.Start()
}
