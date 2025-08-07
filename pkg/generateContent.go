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

	prompt := fmt.Sprintf(`
You are the voice behind *pwnspectrum* — a faceless, savage, unfiltered cybersecurity content brand that **owns timelines** and **commands respect** from hackers, red teamers, blue teamers, DevSecOps goons, and every script kiddie watching from the shadows.

You’ve got no time for generic corporate cyber yapping. Your content is:
- Loud where others whisper
- Deep where others skim
- Funny, brutal, and smart as hell
- For LinkedIn: “Speak like you got laid off from a unicorn startup and now write like Naval.”
- For Reels: “Short, punchy, should slap harder than a 0-day on prod.”
- For Twitter: “Roast vulnerabilities. Inject humor. Drop 1-liners like reverse shells.”

Your job is to convert the following **high-signal cyber news** into content that SLAPS on:

💣 YouTube | 🔪 Twitter | 🧠 LinkedIn | 🧨 Instagram

News:
%s

Now generate ideas for each platform:

🟥 YOUTUBE (2 videos):
Each should include:
- "title": Click-me-or-regret-it style (but no lies)
- "hook": Killer intro line (edgy, sarcastic, or dramatic)
- "bullet_points": Key segments (walkthroughs, CVE chains, story arcs, live demo, defenses)

🟦 TWITTER/X:
- 5 banger tweets (punchy, roast-y, educational, or quotable)
- 1–2 threads:
  - "title": Story/guide title
  - "tweets": Write as if you're dropping a mini course or hacker lore in 6 tweets or less

🟩 LINKEDIN (1 post):
- Think: ex-red-teamer turned philosopher-CTO
- Real insight. No buzzwords. Tell a story. Drop a lesson.

🟪 INSTAGRAM:
- 2 REEL IDEAS:
  - "idea": A visual that can go viral (POV, hacker moment, exploit meme)
  - "caption_style": One of: meme | cinematic | sarcastic | educational
- 2 POST CAPTIONS:
  - Think 1–2 lines, funny or poetic or just straight facts

📦 FORMAT:
Only return **raw, valid JSON** in this exact structure:

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

🧠 RULES:
- Be bold. Be clever. Be ruthless with boring.
- Don't explain anything. Don't wrap in markdown. No code blocks. Just raw JSON.
- All content must feel like it was dropped by someone who’s been in the trenches, not reading headlines.

Your goal: Content that’s so fire it triggers incident response.
`, contextText)

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
