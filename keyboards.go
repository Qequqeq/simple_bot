package main

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func sendSetupKeyboard(bot *telego.Bot, chatID int64, message string) {
	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("✅ Завершить настройку").WithCallbackData("finish_setup"),
		),
	)
	bot.SendMessage(tu.Message(tu.ID(chatID), message).WithReplyMarkup(keyboard))
}

func sendMainKeyboard(bot *telego.Bot, chatID int64, message string) {
	user := userData[chatID]
	var rows [][]telego.InlineKeyboardButton

	if user != nil {
		for _, cat := range user.Categories {
			rows = append(rows, tu.InlineKeyboardRow(
				tu.InlineKeyboardButton(cat).WithCallbackData("category_"+cat),
			))
		}
	}

	rows = append(rows,
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("➕ Добавить категорию").WithCallbackData("add_category"),
			tu.InlineKeyboardButton("📊 Статистика").WithCallbackData("stats"),
		),
	)

	keyboard := tu.InlineKeyboard(rows...)
	bot.SendMessage(tu.Message(tu.ID(chatID), message).WithReplyMarkup(keyboard))
}
