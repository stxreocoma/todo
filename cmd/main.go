package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stxreocoma/todo/auth"
	"github.com/stxreocoma/todo/database"
	"github.com/stxreocoma/todo/handlers"
)

func main() {
	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		port = "7540"
	}

	database.ConnectGorm()

	app := fiber.New()
	app.Static("/", "./web")

	app.Get("/api/nextdate", handlers.GetDate)
	app.Post("/api/task", auth.Authentication(handlers.PostTask))
	app.Get("/api/tasks", auth.Authentication(handlers.GetTasks))
	app.Get("/api/task", auth.Authentication(handlers.GetTask))
	app.Put("/api/task", auth.Authentication(handlers.UpdateTask))
	app.Post("api/task/done", auth.Authentication(handlers.DoneTask))
	app.Delete("api/task", auth.Authentication(handlers.DeleteTask))
	app.Post("api/signin", auth.Registration)

	app.Listen(":" + port)
}
