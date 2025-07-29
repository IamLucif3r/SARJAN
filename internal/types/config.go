package types

type Config struct {
	DatabaseURL       string `json:"db_url"`
	OllamaURL         string `json:"ollama_url"`
	DiscordWebhookURL string `json:"discord_webhook_url"`
	OllamaModel       string `json:"ollama_model"`
}
