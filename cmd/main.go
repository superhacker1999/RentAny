package main

import (
	"RentAny/model/web"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
)

func main() {
	connStr := "postgres://postgres:1234@localhost:5432/mydb12?sslmode=disable"

	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	web.Run()
}
