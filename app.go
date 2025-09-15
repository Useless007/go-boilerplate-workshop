package main

import (
	"boilerplate/database"
	"boilerplate/handlers"
	"boilerplate/middleware"
	"strings"

	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
	prod = flag.Bool("prod", false, "Enable prefork in Production")
)

var skipPaths = []string{
	"/favicon.ico",
	"/img/logo.svg",
	"/.well-known/appspecific/com.chrome.devtools.json",
}

func ShouldSkip(path string) bool {
	for _, p := range skipPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func main() {
	// Parse command-line flags
	flag.Parse()

	// Connected with database
	database.Connect()

	// Create fiber app
	app := fiber.New(fiber.Config{
		Prefork: *prod, // go run app.go -prod
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Next: func(c *fiber.Ctx) bool {
			return ShouldSkip(c.Path())
		},
	}))

	// Create a /api/v1 endpoint
	v1 := app.Group("/api/v1")

	// Bind handlers
	v1.Post("/register", handlers.Register)
	v1.Post("/login", handlers.Login)
	v1.Use(middleware.JWTMiddleware()) // ใช้ middleware ตรวจสอบ JWT token สำหรับทุก endpoint ด้านล่างนี้
	v1.Get("/users", handlers.UserList)
	v1.Get("/me", handlers.Me)
	v1.Post("/users", middleware.RoleMiddleware("ADMIN"), handlers.UserCreate)
	v1.Put("/users/:id", handlers.UserUpdate)
	v1.Delete("/users/:id", handlers.UserDelete)

	// Setup static files
	app.Static("/", "./static/public")

	// Handle not founds
	app.Use(handlers.NotFound)

	// Listen on port 3000
	log.Fatal(app.Listen(*port)) // go run app.go -port=:3000
}
