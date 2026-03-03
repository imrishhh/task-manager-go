package app

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/nullrish/task-manager-go/internal/db"
	"github.com/nullrish/task-manager-go/internal/router"
)

type Server struct {
	IP   string
	Port string
}

func InitializeServer(ip, port string) *Server {
	return &Server{IP: ip, Port: port}
}

func (s *Server) StartServer() {
	addr := fmt.Sprintf("%s:%s", s.IP, s.Port)

	listenConfig := fiber.ListenConfig{
		DisableStartupMessage: true,
		EnablePrefork:         false,
	}

	app := fiber.New(fiber.Config{
		AppName:     "Task Manager Go",
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Khello World!")
	})

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	db := db.InitializeDatabase(connString)

	log.Println("🚡 Configuring Routes....")
	router.ConfigureRoutes(app, db)
	log.Println("Successfully configured routes... ✅")

	fmt.Println()

	log.Println("🚀 Starting Server....")
	log.Printf("Mode: %s\n", app.Config().AppName)
	log.Printf("Address: http://%s\n", addr)
	log.Printf("Prefork: %v\n", listenConfig.EnablePrefork)
	log.Printf("PID: %d\n", os.Getpid())

	log.Println(app.Listen(addr, listenConfig).Error())

	if err := app.Shutdown(); err != nil {
		log.Println("❌ Failed to shutdown app:", err.Error())
	} else {
		log.Println("✅ Successfully shutdowned fiber app.")
	}

	if err := db.Close(); err != nil {
		log.Println("❌ Failed to close database:", err.Error())
	} else {
		log.Println("✅ Successfully closed database pool.")
	}
}
