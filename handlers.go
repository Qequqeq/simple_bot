package main

import (
	"fmt"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func setupHandlers(bh *th.BotHandler) {
	bh.Handle(startHandler, th.CommandEqual("start"))
	bh.Handle(callbackHandler, th.AnyCallbackQuery())
	bh.Handle(messageHandler)
}

func startHandler(bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ID

	userStates[chatID] = &UserState{Mode: "none"}
	userData[chatID] = &UserData{
		Expenses:    make(map[string]float64),
		Categories:  []string{},
		UsingCustom: false,
	}

	keyboard := tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🆓 Кастомные настройки").WithCallbackData("free_mode"),
			tu.InlineKeyboardButton("📦 Базовые категории").WithCallbackData("basic_mode"),
		),
	)

	bot.SendMessage(tu.Message(
		tu.ID(chatID),
		"👋Приветствую Вас! \n\n 🎮Выберите режим работы:",
	).WithReplyMarkup(keyboard))
}

func callbackHandler(bot *telego.Bot, update telego.Update) {
	query := update.CallbackQuery
	if query == nil {
		return
	}

	chatID := query.Message.GetChat().ID
	data := query.Data

	_ = bot.AnswerCallbackQuery(&telego.AnswerCallbackQueryParams{
		CallbackQueryID: query.ID,
	})

	state := userStates[chatID]
	user := userData[chatID]

	if state == nil || user == nil {
		bot.SendMessage(tu.Message(tu.ID(chatID), "⚠️ Пожалуйста, начните с команды /start"))
		return
	}

	switch data {
	case "free_mode":
		handleFreeMode(bot, chatID)
	case "basic_mode":
		handleBasicMode(bot, chatID)
	case "finish_setup":
		handleFinishSetup(bot, chatID)
	case "add_category":
		handleAddCategory(bot, chatID)
	case "stats":
		showStats(bot, chatID)
	default:
		if strings.HasPrefix(data, "category_") {
			handleCategorySelection(bot, chatID, data)
		}
	}
}

func messageHandler(bot *telego.Bot, update telego.Update) {
	msg := update.Message
	if msg == nil {
		return
	}

	chatID := msg.Chat.ID
	state := userStates[chatID]
	user := userData[chatID]

	if state == nil || user == nil {
		bot.SendMessage(tu.Message(tu.ID(chatID), "⚠️ Пожалуйста, начните с команды /start"))
		return
	}

	if state.IsAddingAmount {
		handleAmountInput(bot, chatID, msg.Text)
		state.IsAddingAmount = false
		return
	}

	if state.Mode == "setup" && msg.Text != "" {
		addCategory(chatID, msg.Text)
		sendSetupKeyboard(bot, chatID, "✅ Категория добавлена: "+msg.Text)
	}
}

func handleFreeMode(bot *telego.Bot, chatID int64) {
	userStates[chatID].Mode = "setup"
	userData[chatID].UsingCustom = true
	sendSetupKeyboard(bot, chatID, "✏️ Вводите названия категорий. \n\n🆗Когда закончите, нажмите 'Завершить настройку'. \n\n⚠️Пожалуйста, отнеситесь к настройке внимательно. \n\n🔒После добавления категорий удалить ее будет нельзя!")
}

func handleBasicMode(bot *telego.Bot, chatID int64) {
	userStates[chatID].Mode = "basic"
	userData[chatID].Categories = append([]string{}, defaultCategories...)
	sendMainKeyboard(bot, chatID, "✅ Выбран режим базовых категорий")
}

func handleFinishSetup(bot *telego.Bot, chatID int64) {
	if len(userData[chatID].Categories) == 0 {
		bot.SendMessage(tu.Message(tu.ID(chatID), "❌ Вы должны добавить хотя бы одну категорию!"))
		return
	}

	userStates[chatID].Mode = "free"
	sendMainKeyboard(bot, chatID, "✅ Настройка завершена! Теперь вы можете добавлять расходы")
}

func handleAddCategory(bot *telego.Bot, chatID int64) {
	userStates[chatID].Mode = "setup"
	bot.SendMessage(tu.Message(tu.ID(chatID), "✏️ Введите название новой категории:"))
}

func handleCategorySelection(bot *telego.Bot, chatID int64, data string) {
	category := strings.TrimPrefix(data, "category_")
	userStates[chatID].TempCategory = category
	userStates[chatID].IsAddingAmount = true
	bot.SendMessage(tu.Message(
		tu.ID(chatID),
		fmt.Sprintf("💰 Введите сумму для категории \"%s\": \n\n (Если Вы ошиблись, просто введите отрицательное число)", category),
	))
}
