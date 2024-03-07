package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func ConnectionDB() *sql.DB {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := mysql.Config{
		User:   os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net:    os.Getenv("DB_NET"),
		Addr:   os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	return db

}
