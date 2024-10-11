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

	// TODO : change call of this method to initialization method
	db, err := database.GetConnectionPool()

	if err != nil {
		panic(err)
	}
	defer db.Close()

	web.Run()
}
