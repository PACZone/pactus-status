package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	RPCNodes []string
	BotToken string
	PriceAPI string
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	cfg := Config{
		RPCNodes: strings.Split(os.Getenv("RPC_NODES"), ","),
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		PriceAPI: os.Getenv("PRICE_END_POINT"),
	}

	return cfg, nil
}
