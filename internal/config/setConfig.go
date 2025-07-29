package config

import (
	"fmt"

	"os"

	"github.com/iamlucif3r/sarjan/internal/types"
	"github.com/joho/godotenv"
)

func SetConfig(Config *types.Config) error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("error loading [.env] file: %v", err)
	}
	Config.DatabaseURL = os.Getenv("DATABASE_URL")
	Config.OllamaURL = os.Getenv("OLLAMA_URL")
	Config.DiscordWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL")
	return nil
}
