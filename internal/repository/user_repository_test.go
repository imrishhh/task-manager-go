package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func Test_userRepository(t *testing.T) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load env")
		t.FailNow()
	}
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_NAME"))
	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Printf("Failed to open db pooling")
		t.FailNow()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Printf("Error while pinging database: %s", err)
		log.Printf("Please insure database authentication are correct.")
		log.Printf("Failed to ping database.")
		t.FailNow()
	}
}
