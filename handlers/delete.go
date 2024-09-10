package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stxreocoma/todo/database"
	"github.com/stxreocoma/todo/models"
	"github.com/stxreocoma/todo/utils"
)

func DoneTask(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	var task models.Task

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"error": err.Error()})
	}

	result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ? LIMIT 1", id).First(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(map[string]any{"error": result.Error.Error()})
	}

	if task.Repeat == "" {
		result = database.Gorm.Db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"error": result.Error.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(map[string]any{})

	} else {
		newDate, err := utils.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"error": err.Error()})
		}

		result = database.Gorm.Db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, id)
		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"error": result.Error.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{})
}

func DeleteTask(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"error": err.Error()})
	}

	result := database.Gorm.Db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"error": result.Error.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{})
}
