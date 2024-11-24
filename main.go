package main

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	user "github.com/kyken/fiber-sandbox/handler"
	db "github.com/kyken/fiber-sandbox/lib/db"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// エラーレスポンス用の構造体
type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

var (
	connDB *db.Database
)

func main() {
	config := fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Prefork:     true,
	}
	app := fiber.New(config)

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(cache.New(cache.Config{
		Expiration:   1 * time.Hour,
		CacheControl: true,
	}))

	// CORS
	app.Use(cors.New())

	// Logging Request ID
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	connDB, err := db.NewDatabase()
	if err != nil {
		fmt.Println("failed to connect database: %w", err)
	}

	userHandler := user.NewUserHandler(*connDB)
	app.Get("/users", userHandler.GetUsersHandler)
	app.Get("/user/:id", userHandler.GetUserHandler)
	app.Put("/user", userHandler.PutUserHandler)
	app.Post("/user/:id", userHandler.PostUserHandler)
	app.Delete("/user/:id", userHandler.DeleteUserHandler)

	app.Listen(":3000")
}
