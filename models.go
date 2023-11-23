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
	ChatId    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type Chat struct {
	Id int `json:"id"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Holiday struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Date    string `json:"date"`
	WeekDay string `json:"week_day"`
	Type    string `json:"type"`
}
