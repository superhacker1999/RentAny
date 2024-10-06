package main

import (
	"RentAny/model/database"
	"RentAny/model/web"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file : %v", err)
	}

	// initializes connection pool
	db, err := database.GetConnectionPool()

	if err != nil {
		panic(err)
	}
	defer db.Close()

	web.Run()
}
