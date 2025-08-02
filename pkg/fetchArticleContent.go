package pkg

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/iamlucif3r/sarjan/internal/types"
)

func EnrichArticleContent(articles []types.JudgedArticle) []types.JudgedArticle {
	for i, article := range articles {
		if article.URL == "" {
			log.Printf("Skipping article %s due to empty URL", article.Title)
			continue
		}

		fullContent, err := FetchFullContent(article.URL)
		if err != nil {
			log.Printf("Failed to fetch full content for %s: %v", article.Title, err)
			continue
		}

		articles[i].Content = CleanHTMLContent(fullContent)
	}
	return articles
}

func FetchFullContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func CleanHTMLContent(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html // fallback
	}

	return strings.TrimSpace(doc.Text())
}
