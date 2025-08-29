# Copilot Instructions for SARJAN & TRIKAL

## Big Picture Architecture
- **SARJAN** and **TRIKAL** are Go-based microservices for cybersecurity news curation and notification.
- **TRIKAL**: Fetches, filters, scores, and notifies top cybersecurity news from RSS feeds. Uses local LLM (Ollama) for relevance scoring. Sends Discord notifications.
- **SARJAN**: Focuses on content generation and actionable narratives from curated news. Integrates with TRIKAL via API and Docker Compose.
- Both use PostgreSQL for persistence (`internal/database`), with config in `internal/config`.

## Key Workflows
- **Build SARJAN**: `make build` (see `Makefile`). Binary: `sarjan`.
- **Run SARJAN**: `make run` or `./sarjan`.
- **Docker Compose**: Use `docker-compose.yaml` to run both services together. SARJAN depends on TRIKAL and expects TRIKAL's API at `http://localhost:8089`.
- **TRIKAL**: No Makefile build; run with `go run cmd/main.go` or build manually. Dockerfile provided.
- **RSS Feeds**: Configure in `TRIKAL/rss.yaml`.
- **Discord Alerts**: Set webhook URL in config/env. Alerts sent via `pkg/sendAlerts.go`.
- **LLM Integration**: Ollama API endpoint set via env/config. See `pkg/queryOllama.go`.

## Project-Specific Patterns
- **Config**: Centralized in `internal/config/setConfig.go`. Use `types.Config` struct.
- **Database**: Connection via `internal/database/database.go`. Global `DB` variable.
- **Types**: Shared models in `internal/types/` (e.g., `NewsItem`, `DiscordEmbed`).
- **Batching**: Discord notifications are batched (max 10 embeds per message).
- **Scoring**: LLM scoring logic in `pkg/relevanceScoring.go` and `pkg/queryOllama.go`.
- **Keywords**: Cybersecurity keywords for filtering in `pkg/fetchRSSNews.go`.

## Integration Points
- **Ollama**: Local LLM for scoring, called via HTTP (`pkg/queryOllama.go`).
- **Discord**: Alerts sent via webhook (`pkg/sendAlerts.go`).
- **Docker Compose**: Defines service boundaries and environment variables.

## Conventions & Examples
- **Go Modules**: Each service has its own `go.mod`.
- **Entrypoints**: SARJAN (`cmd/sarjan/main.go`), TRIKAL (`cmd/main.go`).
- **Environment**: Use `.env` files and Docker Compose for secrets/config.
- **Logging**: Use Go's `log` package for service logs.
- **Error Handling**: Return errors up, log failures, continue on batch errors (see Discord alert batching).

## Quick Start
1. Configure RSS feeds (`TRIKAL/rss.yaml`) and Discord webhook.
2. Build SARJAN: `make build`.
3. Run both via Docker Compose: `docker-compose up --build`.
4. Check logs for health and alert status.

---

If any section is unclear or missing, please specify which part needs more detail or examples.
