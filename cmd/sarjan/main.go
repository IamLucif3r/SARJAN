package main

import (
	"context"
	"database/sql"
	"fmt"

	"log"

	"github.com/iamlucif3r/sarjan/internal/config"
	"github.com/iamlucif3r/sarjan/internal/database"
	"github.com/iamlucif3r/sarjan/internal/types"
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

	db := database.DB

	articles, err := pkg.FetchTopRankedArticles(db, 1)
	if err != nil {
		log.Fatalf("Failed to fetch articles: %v", err)
	}

	ctx := context.Background()
	// Generate content ideas
	bundle, err := pkg.GenerateContentIdeas(ctx, articles, *Config)
	if err != nil {
		log.Println("Failed to generate content ideas:", err)

	}
	// ToDo: Add a second model for Viral Check.
	fmt.Println(bundle)

}
