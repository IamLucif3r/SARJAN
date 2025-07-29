package pkg

import (
	"database/sql"

	"github.com/iamlucif3r/sarjan/internal/types"
)

func FetchTopRankedArticles(db *sql.DB, limit int) ([]types.Article, error) {
	query := `SELECT id, title, content, url FROM articles ORDER BY rank DESC LIMIT $1`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []types.Article
	for rows.Next() {
		var art types.Article
		if err := rows.Scan(&art.ID, &art.Title, &art.Content, &art.URL); err != nil {
			return nil, err
		}
		articles = append(articles, art)
	}
	return articles, nil
}
