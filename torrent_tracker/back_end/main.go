package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/amha-mersha/GoTorrent/db"
	"github.com/amha-mersha/GoTorrent/handlers"
	"github.com/amha-mersha/GoTorrent/route"
	"github.com/amha-mersha/GoTorrent/services"
)

func main() {
	if err := LoadEnv(".env"); err != nil {
		log.Fatal("Failed to load environment variables")
	}
	postgres_url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
	Redis, err := db.InitRedis(strings.Join([]string{REIDS_HOST, REDIS_PORT}, ":"), REDIS_PASSWORD, REDIS_DB)
	if err != nil {
		log.Fatal(err)
	}

	PostgreSQL, err := db.InitPostgreSQL(postgres_url)
	if err != nil {
		log.Fatal("Failed to initialize Postgres")
	}

	service := services.NewService(PostgreSQL, Redis)
	handler := handlers.NewHandlers(service)
	r := route.SetupRouter(handler)
	log.Fatal(r.Run(":" + APP_PORT))
}
