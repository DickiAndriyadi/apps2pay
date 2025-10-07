package main

import (
	"apps2pay/handlers"
	"apps2pay/middleware"
	"apps2pay/worker"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgres://myuser:password@localhost:5432/cinema_db?sslmode=disable")
	if err != nil {
		log.Fatal("❌ Cannot connect to DB:", err)
	}
	defer conn.Close(context.Background())
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	handlers.SetDB(conn, redisClient)

	app := fiber.New()

	go worker.StartSeatCleanupWorker(context.Background())

	// Public route
	app.Post("/login", handlers.Login)

	// Protected routes
	api := app.Group("/api")
	api.Use(middleware.AuthMiddleware())

	// Public untuk semua user terautentikasi
	api.Get("/schedules", handlers.GetSchedules)
	api.Get("/schedules/:id", handlers.GetSchedule)
	api.Post("/schedules/:id/seats/lock", handlers.LockSeat)
	api.Post("/schedules/:id/seats/release", handlers.ReleaseSeat)
	api.Get("/schedules/:id/seats", handlers.GetSeatsBySchedule)
	api.Post("/schedules/:id/seats/confirm", handlers.ConfirmSeat)

	// Hanya admin
	admin := api.Group("")
	admin.Use(middleware.AdminOnly())
	admin.Post("/schedules", handlers.CreateSchedule)
	admin.Put("/schedules/:id", handlers.UpdateSchedule)
	admin.Delete("/schedules/:id", handlers.DeleteSchedule)
	admin.Post("/schedules/:id/cancel", handlers.CancelSchedule)

	log.Println("🚀 Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))

}
