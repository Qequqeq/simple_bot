package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func addCategory(chatID int64, name string) {
	user := userData[chatID]
	if user == nil {
		return
	}

	for _, cat := range user.Categories {
		if cat == name {
			return
		}
	}
	user.Categories = append(user.Categories, name)
}

func handleAmountInput(bot *telego.Bot, chatID int64, text string) {
	state := userStates[chatID]
	user := userData[chatID]

	if state == nil || user == nil {
		return
	}

	sum, err := strconv.ParseFloat(text, 64)
	if err != nil {
		bot.SendMessage(tu.Message(tu.ID(chatID), "❌ Некорректная сумма!"))
		return
	}

	user.Expenses[state.TempCategory] += sum
	sendMainKeyboard(bot, chatID, fmt.Sprintf("✅ Добавлено %.2f ₽ в категорию \"%s\" \n\n⚠️Не забывайте выбрать новую категорию!", sum, state.TempCategory))
}

func showStats(bot *telego.Bot, chatID int64) {
	user := userData[chatID]
	if user == nil {
		return
	}

	var builder strings.Builder
	builder.WriteString("📊 Статистика:\n\n")

	total := 0.0
	for _, cat := range user.Categories {
		sum := user.Expenses[cat]
		builder.WriteString(fmt.Sprintf("▫️ %s: %.2f ₽\n", cat, sum))
		total += sum
	}

	builder.WriteString(fmt.Sprintf("\n💵 Итого: %.2f ₽", total))
	sendMainKeyboard(bot, chatID, builder.String())
}
