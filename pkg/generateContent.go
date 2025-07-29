package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/iamlucif3r/sarjan/internal/types"
)

func GenerateContentIdeas(ctx context.Context, articles []types.Article, Config types.Config) (types.ContentIdeas, error) {
	var contextText string
	for _, article := range articles {
		contextText += fmt.Sprintf("- %s: %s\n", article.Title, article.Content)
	}

	prompt := fmt.Sprintf(`Given the following cybersecurity news headlines and summaries:

	%s

	Generate:
	1. A detailed YouTube video idea with a hook, title, and bullet points.
	2. Three short and punchy Tweets to share.
	3. A detailed LinkedIn post draft.
	4. Two Instagram Reel ideas and caption styles.

	Respond strictly in this JSON format:
	{
	"YouTube": "...",
	"Twitter": "...",
	"LinkedIn": "...",
	"Instagram": "..."
	}
	`, contextText)

	ollamaURL := Config.OllamaURL
	model := Config.OllamaModel

	payload := map[string]interface{}{
		"model":  model,
		"prompt": prompt,
		"stream": false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return types.ContentIdeas{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ollamaURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return types.ContentIdeas{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return types.ContentIdeas{}, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.ContentIdeas{}, fmt.Errorf("failed to read response: %w", err)
	}

	type ollamaResponse struct {
		Response string `json:"response"`
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return types.ContentIdeas{}, fmt.Errorf("failed to parse Ollama response: %w", err)
	}

	var ideas types.ContentIdeas
	if err := json.Unmarshal([]byte(ollamaResp.Response), &ideas); err != nil {
		return types.ContentIdeas{}, fmt.Errorf("failed to parse model JSON output: %w", err)
	}

	return ideas, nil
}
