package api

type Update struct {
	UpdateId int           `json:"update_id"`
	Message  IncomeMessage `json:"message"`
}

type IncomeMessage struct {
	Text     string   `json:"text"`
	Chat     Chat     `json:"chat"`
	Date     int      `json:"date"`
	Location Location `json:"location"`
}
type OutgoingMessage struct {
	ChatId    int      `json:"chat_id"`
	Text      string   `json:"text"`
	ParseMode string   `json:"parse_mode"`
	Keyboard  Keyboard `json:"reply_markup,omitempty"`
}

type Chat struct {
	Id int `json:"id"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type KeyboardButton struct {
	Text string `json:"text"`
}

type Keyboard struct {
	InlineKeyboard [][]KeyboardButton `json:"keyboard"`
}
