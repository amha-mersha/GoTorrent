package main

import (
	"fmt"
	"log"

	"github.com/amha-mersha/GoTorrent/db"
)

func main() {
	LoadEnv(".env")
	postgres_url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
	if err := db.InitSQLite(); err != nil {
		log.Fatal("Failed to initialize SQLite")
	}
	if err := db.InitPostgreSQL(postgres_url); err != nil {
		log.Fatal("Failed to initialize Postgres")
	}

}
