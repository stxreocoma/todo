package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stxreocoma/todo/database"
	"github.com/stxreocoma/todo/handlers"
)

func main() {
	//database.ConnectDB()
	//defer database.DB.Db.Close()

	database.ConnectGorm()

	app := fiber.New()
	app.Static("/", "./web")

	log.Println(os.Getenv("TODO_PORT"))

	app.Get("/api/nextdate", handlers.Date)
	app.Post("/api/task", handlers.PostTask)
	app.Get("/api/tasks", handlers.GetTasks)
	app.Get("/api/task", handlers.GetTask)
	app.Put("/api/task", handlers.UpdateTask)

	app.Listen(":7540")
}
