package main

import (
	"RentAny/internal/controller"
	"RentAny/internal/repository/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func waitForDBToStartup() (*postgres.Database, error) {
	waitTimeout, err := strconv.Atoi(os.Getenv("RENT_ANY_DB_STARTUP_WAIT_TIMEOUT"))

	if err != nil {
		log.Println("Could not parse \"RENT_ANY_DB_STARTUP_WAIT_TIMEOUT\", default parameter 30 seconds was set")
		waitTimeout = 30
	}

	start := time.Now()

	for {
		db, err := postgres.GetConnectionPool()

		if err != nil {
			if time.Since(start) >= (time.Duration(waitTimeout) * time.Second) {
				return nil, err
			}
		} else {
			return db, nil
		}
	}

}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file : %v", err)
	}

	// TODO : change call of this method to initialization method
	db, err := waitForDBToStartup()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server.Run()
}
