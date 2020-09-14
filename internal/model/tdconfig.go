package model

// TDConfig tdlib configuration
type TDConfig struct {
	TelegramAPIID   string `toml:"telegram_api_id"`
	TelegramAPIHash string `toml:"telegram_api_hash"`

	PhoneNumber  string `toml:"telegram_phone_number"`
	AuthPassword string `toml:"telegram_password"`
	AuthCode     string
}
