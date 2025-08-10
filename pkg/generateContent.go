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
- For LinkedIn: “Speak like you got laid off from a unicorn startup and now write like Naval.” Grounded in **practical, operational tactics** — share methods, frameworks, and war stories tied to the news.
- For Reels: “Short, punchy, should slap harder than a 0-day on prod.” Visually engaging, instantly digestible, but hinting at a bigger play hackers will want to dig into.
- For Twitter: “Roast vulnerabilities. Inject humor. Drop 1-liners like reverse shells.” Tactical, witty, and dripping with hacker culture references.

You will:
- Extract not just *what happened*, but the **real operational impact** for attackers and defenders.
- Call out possible exploitation paths, detection methods, mitigation tips, or counter-tactics — while staying platform-specific.
- Offer **multiple interpretations** of the same news so each platform gets its own angle.

🔍 **Before writing anything, run this mental checklist**:
1. **Attack Chain** — How could this be exploited end-to-end? What steps would an attacker take? What tooling or TTPs fit here?
2. **Detection Gap** — How would most defenders miss this? Where are logging, monitoring, or response weaknesses?
3. **Mitigation** — How could an org patch, detect, or harden against it *now* without waiting for a vendor fix?
4. Translate those insights into platform-specific content without explicitly stating the checklist.

Your job is to convert the following **high-signal cyber news** into content that SLAPS on:

💣 YouTube | 🔪 Twitter | 🧠 LinkedIn | 🧨 Instagram

News:
%s

Now generate ideas for each platform:

🟥 YOUTUBE (2 videos):
Each should include:
- "title": Click-me-or-regret-it style (but no lies)
- "hook": Killer intro line (edgy, sarcastic, or dramatic) that teases the tactical angle
- "bullet_points": Story beats showing exploitation flow, real-world attack scenarios, or defense breakdowns

🟦 TWITTER/X:
- 5 banger tweets (mix humor + actionable takeaway — e.g., an exploit vector, detection tip, or TTP summary)
- 1–2 threads:
  - "title": Story/guide title with curiosity baked in
  - "tweets": Drop a war story, a condensed exploit walkthrough, or “how to spot/fix” guide in 6 tweets or less — every tweet adds value

🟩 LINKEDIN (1 post):
- Tactical but framed for professionals
- Tell a short, impactful story from the news with a hacker’s lens — highlight the exploitation chain, operational blind spots, and the lesson for defenders

🟪 INSTAGRAM:
- 2 REEL IDEAS:
  - "idea": Visual hook (POV exploit moment, hacker POV, meme-worthy attack chain)
  - "caption_style": meme | cinematic | sarcastic | educational — match to the operational angle
- 2 POST CAPTIONS:
  - 1–2 lines, either savage or surgical — must hit emotionally or technically

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
- Focus on **practical, operational insights** — no vague “awareness” fluff.
- Don’t just summarize — show how attackers would weaponize it, and how defenders can counter.
- No markdown, no explanations, no code blocks — just raw JSON.
- Content must read like it came from someone who lives in exploits, packets, and logs — not news headlines.

Your goal: Content so tactical and savage it gets bookmarked by pentesters, banned in corporate Slack, and screenshot into threat intel decks without credit.
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
