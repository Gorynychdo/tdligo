package model

type Config struct {
	TelegramAPIID   string `toml:"telegram_api_id"`
	TelegramAPIHash string `toml:"telegram_api_hash"`

	PhoneNumber  string `toml:"telegram_phone_number"`
	AuthPassword string `toml:"telegram_password"`
	AuthCode     string

	HTTPPort string `toml:"http_port"`
}

func NewConfig() *Config {
	return &Config{
		HTTPPort: ":8000",
	}
}
