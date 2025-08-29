package types

type ContentIdeas struct {
	InstagramReels []struct {
		Idea         string `json:"idea"`
		CaptionStyle string `json:"caption_style"`
	} `json:"instagram_reels"`

	InstagramPosts []string `json:"instagram_posts"`

	TwitterPosts []string `json:"twitter_posts"`

	TwitterThreads []struct {
		Title string   `json:"title"`
		Body  []string `json:"body"`
	} `json:"twitter_threads"`

	LinkedInPosts []string `json:"linkedin_posts"`

	YouTubeVideoIdeas []struct {
		Title        string   `json:"title"`
		Hook         string   `json:"hook"`
		BulletPoints []string `json:"bullet_points"`
	} `json:"youtube_video_ideas"`
}
