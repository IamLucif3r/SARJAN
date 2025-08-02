package main

import (
	"context"
	"database/sql"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/iamlucif3r/sarjan/internal/config"
	"github.com/iamlucif3r/sarjan/internal/database"
	"github.com/iamlucif3r/sarjan/internal/types"
	"github.com/iamlucif3r/sarjan/internal/utils"
	"github.com/iamlucif3r/sarjan/pkg"
)

var Config *types.Config
var Db *sql.DB

func init() {
	log.Println("Initializing SARJAN...")
	log.Println("Initializing configuration...")
	Config = &types.Config{}

	err := config.SetConfig(Config)
	if err != nil {
		log.Printf("Error setting configuration: %v\n", err)
		return
	}
	Db, err = database.ConnectDB(*Config)
	if err != nil {
		log.Printf("Error connecting to database: %v\n", err)
		return
	}
	log.Println("Configuration initialized successfully.")
}

func main() {

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "SARJAN : Smart Assistant for Real-time Journey from news to Actionable Narratives",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Alive",
		})
	})

	router.POST("/generate", func(c *gin.Context) {

		log.Println("[INFO] Starting content generation process...")

		db := database.DB
		ctx := context.Background()

		articles, err := pkg.FetchTopRankedArticles(db, 1)
		if err != nil {
			log.Fatalf("Failed to fetch articles: %v", err)
		}
		log.Println("[INFO] Fetched ", len(articles), " articles from database")

		bundle, err := pkg.GenerateContentIdeas(ctx, articles, *Config)
		if err != nil {
			log.Println("[Error] Failed to generate content ideas:", err)

		}
		log.Println("[INFO] Generated content ideas successfully")
		err = utils.GenerateContentIdeasPDF(bundle, "output/content_ideas.pdf")
		if err != nil {
			log.Println("Failed to generate PDF:", err)
		}

		err = utils.SendPDFToDiscord(Config.DiscordWebhookURL, "output/content_ideas.pdf")
		if err != nil {
			log.Println("Failed to send Markdown to Discord:", err)
		} else {
			log.Println("[Info] Sent content ideas to Discord successfully!")
		}

	})
	gin.SetMode(gin.ReleaseMode)
	router.Run(":4444")
}
