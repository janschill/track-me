package main

import (
	"log"

	"github.com/janschill/track-me/internal/config"
	"github.com/janschill/track-me/internal/db"
	"github.com/janschill/track-me/internal/repository"
	"github.com/janschill/track-me/internal/service"
	_ "github.com/mattn/go-sqlite3"
)

var conf *config.Config

func init() {
	var err error
	conf, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Couldnt load config %v", err)
	}
}

func main() {
	if conf.DatabaseURL == "" {
		log.Fatal("DB_PATH environment variable is not set")
	}
	Db, err := db.InitializeDB(conf.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	day := "2024-08-05"

	repo := repository.NewRepository(Db)
	aggregationService := service.NewAggregationService(repo)

	aggregationService.Aggregate(day)
}
