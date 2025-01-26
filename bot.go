package main

import (
	"fmt"
	"os"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	botToken := "<your_bot_token>"
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	updates, _ := bot.UpdatesViaLongPolling(nil)
	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
	}, th.CommandEqual("start"))

	for update := range updates {
		if update.Message != nil {
			chatID := tu.ID(update.Message.Chat.ID)

			_, _ = bot.SendSticker(
				tu.Sticker(
					chatID,
					tu.FileFromID("CAACAgIAAxkBAAENo39nle9tT-QnxDNLv_Bjt7_C_QdznwACCwADLfpWFAlJMEy8OeyNNgQ"),
				),
			)
		}
	}
	bh.Start()

}

