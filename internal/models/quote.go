package models

// Quote - структура для хранения цитат
type Quote struct {
	ID     int    `json:"id"`
	Author string `json:"author"`
	Text   string `json:"quote"`
}
