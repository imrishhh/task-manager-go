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
	"github.com/nullrish/task-manager-go/internal/model"
	"github.com/nullrish/task-manager-go/internal/repository"
)

func Test_userRepository(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("Failed to load env: %v", err)
		t.FailNow()
	}
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_NAME"))
	db, err := sql.Open("pgx", connString)
	defer db.Close()
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

	userRepo := repository.NewUserRepository(db)

	req := model.UserRequest{
		Username: "tester1",
		Password: "tester244#O",
		Email:    "mail@mail.com",
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if user, err := userRepo.CreateUser(ctx, &req); err != nil {
		log.Printf("Failed to create user: %v", err)
		t.FailNow()
	} else {
		log.Println("[CreateUser] executed successfully!", user)
	}

	user, err := userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		t.FailNow()
	}
	log.Println("[GetUserByUsername] executed successfully!", user)

	user, err = userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("Failed to get user by email: %v", err)
		t.FailNow()
	}
	log.Println("[GetUserByEmail] executed successfully!", user)

	req = model.UserRequest{
		Username: "tester2",
		Password: "fakfak",
		Email:    "mail.mail@mail.com",
	}

	user, err = userRepo.UpdateUser(ctx, user.ID, &req)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		t.FailNow()
	}
	log.Println("[UpdateUser] executed successfully!", user)

	err = userRepo.DeleteUser(ctx, user.ID)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		t.FailNow()
	}
	log.Println("[DeleteUser] executed successfully!")
}
