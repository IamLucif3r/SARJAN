package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/iamlucif3r/sarjan/internal/types"
)

func JudgeArticlesComparatively(articles []types.Article, modelName string) ([]types.JudgedArticle, error) {
	scored := make([]types.JudgedArticle, len(articles))
	for i := range articles {
		scored[i] = types.JudgedArticle{Article: articles[i], Score: 0}
	}

	articleJSON, err := json.MarshalIndent(articles, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal articles: %w", err)
	}

	prompt := fmt.Sprintf(`
You are a cybersecurity strategist working for the faceless threat intel brand "pwnspectrum".

You are given multiple real-world cybersecurity articles in JSON format. Your job is to evaluate and score them **comparatively** across these criteria:

1. Relevance to current cybersecurity threats  
2. Uniqueness and novelty  
3. Technical depth (exploits, root cause, complexity)  
4. Viral content potential (LinkedIn, YouTube, Twitter)  
5. Actionability for defenders and researchers  
6. Timeliness (emerging or trending issues)

🎯 TASK:
- Score each article **relative to others**, not in isolation.
- Use scores from **1 (weak)** to **10 (strong)**.

📦 RESPONSE FORMAT:
- Respond **only with raw JSON**. No text, no headings, no code block formatting.
- Output **must match** this exact format:

{"Article 1": 7, "Article 2": 9, "Article 3": 6}

🚫 DO NOT:
- Write explanations, comments, markdown, or natural language.
- Wrap the JSON in triple backticks.
- Add quotes or notes outside the JSON object.

Now here are the articles in JSON format:
%s
`, string(articleJSON))

	scoreMap, err := QueryOllamaScoreMap(prompt)
	if err != nil {
		return scored, fmt.Errorf("failed to query Ollama for scoring: %w", err)
	}
	for i := range scored {
		key := fmt.Sprintf("Article %d", i+1)
		if score, ok := scoreMap[key]; ok {
			scored[i].Score = int(score)
		}
	}
	log.Println("[DEBUG] Scored articles before sorting:", scored)
	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	return scored, nil
}
