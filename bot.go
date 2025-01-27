package main

import (
	"fmt"
	"os"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	botToken := "7942804600:AAHPdp_AdFl47YgkMJxaqCfaRGDHEw2mBJo" // <your_bot_token>
	if botToken == "" {
		fmt.Println("BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)
	bh, _ := th.NewBotHandler(bot, updates)

	setupHandlers(bh)

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Start()
}
