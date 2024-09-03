package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stxreocoma/todo/database"
	"github.com/stxreocoma/todo/models"
	"github.com/stxreocoma/todo/validation"
)

func PostTask(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	var task models.Task

	err := c.BodyParser(&task)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": err.Error()})
	}

	task.Date, err = validation.Date(task)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": err.Error()})
	}

	if len(task.Title) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": models.ErrTitle})
	}

	err = validation.Repeat(task.Repeat)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": err.Error()})
	}

	result := database.Gorm.Db.Omit("created_at", "updated_at", "deleted_at").Create(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"Error": result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{"ID": task.ID})
}
