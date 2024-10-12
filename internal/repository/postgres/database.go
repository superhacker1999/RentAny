package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

var singleConnectionPool *sqlx.DB

type Database struct {
	db      *sqlx.DB
	userDAO *UserDAO
}

func (db *Database) initDAOs() {
	db.userDAO = newUserDAO(db.db)

	// add new DAOs here
}

// returns pointer to singleton connection pool
func GetConnectionPool() (*Database, error) {
	database := &Database{}

	if singleConnectionPool != nil {
		database.db = singleConnectionPool
		database.initDAOs()
		return database, nil
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSL_MODE")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, sslMode)

	var err error
	singleConnectionPool, err = sqlx.Connect("pgx", connStr)

	if err != nil {
		return nil, err
	}

	singleConnectionPool.SetMaxOpenConns(25)
	singleConnectionPool.SetMaxIdleConns(25)
	singleConnectionPool.SetConnMaxIdleTime(1 * time.Second)

	database.db = singleConnectionPool
	database.initDAOs()
	return database, nil
}

func (db *Database) Close() {
	db.db.Close()
}

func (db *Database) GetUserDAO() (*UserDAO, error) {
	if db.userDAO != nil {
		return db.userDAO, nil
	}
	return nil, fmt.Errorf("userDAO not initialized")
}
