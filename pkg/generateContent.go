package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/iamlucif3r/sarjan/internal/types"
)

func GenerateContentIdeas(ctx context.Context, articles []types.JudgedArticle, Config types.Config) (types.ContentIdeas, error) {
	var contextText string
	var contentIdea types.ContentIdeas
	for _, article := range articles {
		contextText += fmt.Sprintf("- %s: %s\n", article.Title, article.Content)
	}

	prompt := fmt.Sprintf(`You are a cybersecurity content strategist and social media growth hacker helping a faceless personal brand named "pwnspectrum" generate impactful, engaging, and platform-tailored content ideas.
The audience includes cybersecurity professionals, ethical hackers, researchers, developers, and tech-savvy individuals.

Given the following curated, high-impact cybersecurity news headlines and summaries:

%s

Generate content ideas for these platforms:

1. YouTube:
- 2 Video ideas
- Each must include:
  - A powerful "title" (clickbait-style but not misleading)
  - A "hook" (curiosity-building one-liner)
  - A list of "bullet_points" (3–5 points explaining the video structure)

2. Twitter/X:
- 5 individual tweet ideas as strings (short, witty, or value-packed)
- 1 or 2 tweet threads:
  - Each should be a JSON object with:
    - "title": Title of the thread
    - "tweets": Array of tweet strings (structured as a story or tutorial)

3. LinkedIn:
- 1 complete post (professional tone, thought-leadership style)

4. Instagram:
- 2 Reel ideas:
  - Each should include:
    - "idea": The concept or idea of the reel
    - "caption_style": Type of caption (e.g., funny, sarcastic, insightful)
- 2 Post captions:
  - Each should be a short, 1–2 line Instagram-native caption

Important Rules:
- Use a casual tone for Twitter and Instagram.
- Use a professional and insightful tone for LinkedIn.
- Make YouTube titles punchy, not exaggerated clickbait.
- DO NOT include any explanation, markdown, or extra text. Only return valid JSON in the exact format below.

The expected JSON structure:

{
  "linkedin_posts": ["string"],
  "youtube_video_ideas": [
    {
      "title": "string",
      "hook": "string",
      "bullet_points": ["string", "string", "string"]
    },
    {
      "title": "string",
      "hook": "string",
      "bullet_points": ["string", "string", "string"]
    }
  ],
  "instagram_reels": [
    {
      "idea": "string",
      "caption_style": "string"
    },
    {
      "idea": "string",
      "caption_style": "string"
    }
  ],
  "instagram_posts": ["string", "string"],
  "twitter_posts": ["string", "string", "string", "string", "string"],
  "twitter_threads": [
    {
      "title": "string",
      "tweets": ["string", "string", "string"]
    }
  ]
}

Keep it crisp, high-value, engaging, and tailored for virality.`, contextText)

	ollamaURL := Config.OllamaURL + "/api/generate"

	payload := map[string]interface{}{
		"model":       Config.OllamaModel,
		"prompt":      prompt,
		"temperature": 0.7,
		"stream":      false,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return types.ContentIdeas{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(ollamaURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return contentIdea, fmt.Errorf("failed to call Ollama API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return contentIdea, fmt.Errorf("non-200 response: %d %s", resp.StatusCode, string(bodyBytes))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return contentIdea, fmt.Errorf("failed to read response body: %v", err)
	}

	type ollamaResponse struct {
		Response string `json:"response"`
	}

	// fmt.Println("RAW OLLAMA:", string(respBody))

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return types.ContentIdeas{}, fmt.Errorf("failed to parse Ollama response: %w", err)
	}
	fmt.Println("OLLAMA RESPONSE:", ollamaResp.Response)
	cleanOutput := strings.TrimSpace(ollamaResp.Response)
	if strings.HasPrefix(cleanOutput, "```") {
		cleanOutput = strings.TrimPrefix(cleanOutput, "```json")
		cleanOutput = strings.TrimPrefix(cleanOutput, "```")
		cleanOutput = strings.TrimSuffix(cleanOutput, "```")
		cleanOutput = strings.TrimSpace(cleanOutput)
	}

	fmt.Println("CLEANED OLLAMA OUTPUT:", cleanOutput)
	var tmp map[string]any
	if err := json.Unmarshal([]byte(cleanOutput), &tmp); err != nil {
		log.Fatal("Invalid JSON structure:", err)
	}
	for k := range tmp {
		log.Println("[DEBUG] Key in JSON:", k)
	}

	var ideas types.ContentIdeas
	if err := json.Unmarshal([]byte(cleanOutput), &ideas); err != nil {
		log.Println("[ERROR] Failed to parse output:", err)
		return types.ContentIdeas{}, fmt.Errorf("failed to parse model JSON output: %w", err)
	}

	return ideas, nil
}
