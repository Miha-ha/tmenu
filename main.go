package tmenu

import (
	"gopkg.in/telegram-bot-api.v4"
)

type Menu struct {
	Back     *Menu
	items    [][]*MenuItem
	itemsMap map[string]*MenuItem
}

type ActionParams struct {
	Menu   *Menu
	Item   *MenuItem
	Update tgbotapi.Update
}

type MenuItem struct {
	Text   string
	Data   string
	Tag    int
	Action func(param *ActionParams) *Menu
}

func NewMenu(parent *Menu, items [][]*MenuItem) *Menu {
	menu := &Menu{
		Back: parent,
	}

	menu.SetItems(items)

	return menu
}

func (m *Menu) SetItems(items [][]*MenuItem) {

	m.items = make([][]*MenuItem, len(items))
	itemsMap := map[string]*MenuItem{}

	for i, row := range items {
		m.items[i] = make([]*MenuItem, len(row))
		for j, item := range row {
			if item.Data == "" {
				item.Data = item.Text
			}

			m.items[i][j] = item
			itemsMap[item.Data] = item
		}
	}

	m.itemsMap = itemsMap
}

func (m *Menu) AddSub(items [][]*MenuItem) *Menu {
	return NewMenu(m, items)
}

func (m *Menu) Draw(update tgbotapi.Update) tgbotapi.Chattable {
	if update.CallbackQuery == nil {
		return nil
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(m.items))
	for _, r := range m.items {
		buttons := make([]tgbotapi.InlineKeyboardButton, 0, len(r))
		for _, c := range r {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(c.Text, c.Data))
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(buttons...))

	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	// update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, keyboardInlineSub1
	return tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, keyboard)
}

func (m *Menu) Process(update tgbotapi.Update) *Menu {
	if update.CallbackQuery != nil {
		if item, ok := m.itemsMap[update.CallbackQuery.Data]; ok && item.Action != nil {
			return item.Action(&ActionParams{
				Menu:   m,
				Item:   item,
				Update: update,
			})
		}
	}

	return m
}
