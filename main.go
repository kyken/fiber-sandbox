package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/storage/redis/v3"
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
	// middleware にcache 機構があるが、エンドポイント個別に制御したほうが柔軟なので個別定義とする
	connCache *redis.Storage
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

	// CORS
	app.Use(cors.New())

	// Logging Request ID
	// app.Use(requestid.New())
	// app.Use(logger.New(logger.Config{
	// 	// For more options, see the Config section
	// 	Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	// }))

	fiberlog.SetLevel(fiberlog.LevelDebug)
	app.Use(logger.New(logger.Config{}))

	connDB, err := db.NewDatabase()
	if err != nil {
		fmt.Println("failed to connect database: %w", err)
	}

	connCache = redis.New(redis.Config{
		Host:      "127.0.0.1",
		Port:      6379,
		Username:  "",
		Password:  "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})
    app.Use(cache.New(cache.Config{
        Storage: connCache,
		Expiration: 1 * time.Hour,
		Next: func(c *fiber.Ctx) bool {
			return c.Method() != "GET"
		},
        // キャッシュのキーを生成する関数をカスタマイズ
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.OriginalURL()
        },
    }))

	userHandler := user.NewUserHandler(*connDB, *connCache)
	app.Get("/users", userHandler.GetUsersHandler)
	app.Get("/user/:id", userHandler.GetUserHandler)
	app.Put("/user", userHandler.PutUserHandler)
	app.Post("/user/:id", userHandler.PostUserHandler)
	app.Delete("/user/:id", userHandler.DeleteUserHandler)

	app.Listen(":3000")
}
