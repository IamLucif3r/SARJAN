package types

type Article struct {
	ID      int
	Title   string
	Content string
	URL     string
}

type JudgedArticle struct {
	Article
	Score      int
	FinalScore float64
}
