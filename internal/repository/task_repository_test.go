package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/nullrish/task-manager-go/internal/model"
	"github.com/nullrish/task-manager-go/internal/repository"
)

func Test_taskRepository(t *testing.T) {
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

	taskRepo := repository.NewTaskRepository(db)
	userRepo := repository.NewUserRepository(db)

	userReq := model.UserRequest{
		Username: "tester1",
		Password: "adadada",
		Email:    "MailerMailed@mail.com",
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if _, err := db.Exec("DELETE FROM tasks;"); err != nil {
		log.Printf("%v", err)
		t.FailNow()
	}

	if _, err := db.Exec("DELETE FROM users;"); err != nil {
		log.Printf("%v", err)
		t.FailNow()
	}

	user, err := userRepo.CreateUser(ctx, &userReq)
	if err != nil {
		log.Printf("failed to create user: %v", err)
		t.FailNow()
	}
	log.Println("[CreateUser] executed successfully!", user)
	defer userRepo.DeleteUser(ctx, user.ID)

	req := model.TaskRequest{
		TaskTitle:       "A test title",
		TaskDescription: "A test description",
		Status:          "pending",
		UserID:          user.ID,
	}

	task, err := taskRepo.CreateTask(ctx, &req)
	if err != nil {
		log.Printf("Failed to create task: %v", err)
		t.FailNow()
	}
	log.Println("[CreateTask] executed successfully!", task)

	task, err = taskRepo.GetTaskByID(ctx, task.ID)
	if err != nil {
		log.Printf("Failed to get task by id: %v", err)
		t.FailNow()
	}
	log.Println("[GetTaskByID] executed successfully!", task)

	tasks, err := taskRepo.GetTasks(ctx)
	if err != nil {
		log.Printf("Failed to get tasks: %v", err)
		t.FailNow()
	}
	log.Println("[GetTasks] Fetched all tasks", tasks)

	req.TaskTitle = "New Test Title"
	req.ID = task.ID

	task, err = taskRepo.UpdateTask(ctx, &req)
	if err != nil {
		log.Printf("Failed to update tasks: %v", err)
		t.FailNow()
	}
	log.Println("[UpdateTask] updated task!", task)

	tasks, err = taskRepo.GetTasksByUserID(ctx, user.ID)
	if err != nil {
		log.Printf("Failed to get tasks: %v", err)
		t.FailNow()
	}
	log.Println("[GetTasksByUserID] get tasks!", tasks)

	err = taskRepo.DeleteTask(ctx, task.ID)
	if err != nil {
		log.Println("Failed to delete task")
		t.FailNow()
	}
	log.Println("[DeleteTask] deleted task!")
}
