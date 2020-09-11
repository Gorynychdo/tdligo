package model

// TelegramSettings настройки учетной записи Telegram
type TelegramSettings struct {
	Number   string `json:"number"`
	Code     string `json:"code"`
	Password string `json:"password"`
}
