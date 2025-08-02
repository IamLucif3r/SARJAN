package pkg

import (
	"database/sql"
	"fmt"

	"github.com/iamlucif3r/sarjan/internal/types"
)

func FetchTopRankedArticles(db *sql.DB, limit int) ([]types.JudgedArticle, error) {
	var maxScore float64
	err := db.QueryRow(`SELECT MAX(llm_score) FROM articles`).Scan(&maxScore)
	if err != nil {
		return nil, fmt.Errorf("failed to query max llm_score: %v", err)
	}

	query := `SELECT id, title, description, link FROM articles WHERE llm_score = $1`
	rows, err := db.Query(query, maxScore)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []types.JudgedArticle
	for rows.Next() {
		var art types.JudgedArticle
		if err := rows.Scan(&art.ID, &art.Title, &art.Content, &art.URL); err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}
	return articles, nil
}
