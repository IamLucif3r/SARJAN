package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/iamlucif3r/sarjan/internal/types"
)

var DB *sql.DB

func ConnectDB(cfg types.Config) (*sql.DB, error) {
	connStr := cfg.DatabaseURL
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// err = DB.Ping()
	// if err != nil {
	// 	return nil, fmt.Errorf("error pinging database: %v", err)
	// }

	log.Println("Successfully connected to the database")

	return DB, nil
}
