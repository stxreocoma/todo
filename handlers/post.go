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
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"error": err.Error()})
	}

	task.Date, err = validation.Date(task)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"error": err.Error()})
	}

	if len(task.Title) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"error": models.ErrTitle.Error()})
	}

	err = validation.Repeat(task.Repeat)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"error": err.Error()})
	}

	result := database.Gorm.Db.Omit("created_at", "updated_at", "deleted_at", "id").Create(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"error": result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{"id": task.ID})
}
