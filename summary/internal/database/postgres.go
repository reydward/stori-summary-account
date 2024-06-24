package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func NewPostgresConnection() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	log.Printf("NewPostgresConnection()	db: %v", db)
	if err != nil {
		log.Printf("Error getting the db connection %s", err.Error())
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Error pinging the db %s", err.Error())
		return nil, err
	}

	return db, nil
}
