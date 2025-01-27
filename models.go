package main

type UserState struct {
	Mode           string
	TempCategory   string
	IsAddingAmount bool
}

type UserData struct {
	Categories  []string
	Expenses    map[string]float64
	UsingCustom bool
}

var (
	userStates = make(map[int64]*UserState)
	userData   = make(map[int64]*UserData)

	defaultCategories = []string{
		"дом", "налоги", "аптека", "подарки",
		"продукты", "транспорт", "одежда", "досуг", "другое",
	}
)
