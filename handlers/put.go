package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stxreocoma/todo/database"
	"github.com/stxreocoma/todo/models"
	"github.com/stxreocoma/todo/utils"
	"github.com/stxreocoma/todo/validation"
)

func UpdateTask(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	var task models.Task

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": models.ErrTaskNotFound})
	}

	date, err := time.Parse(utils.DateFormat, task.Date)
	if err != nil {
		if task.Date == "" {
			date = time.Now()
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": err.Error()})
		}
	} else if date.Unix() < time.Now().Unix() {
		if task.Repeat == "" {
			date = time.Now()
		} else {
			dateString, err := utils.NextDate(time.Now(), date.Format(utils.DateFormat), task.Repeat)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": err.Error()})
			}
			date, err = time.Parse(utils.DateFormat, dateString)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": err.Error()})
			}
		}
	}
	task.Date = date.Format(utils.DateFormat)

	if len(task.Title) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": models.ErrTitle})
	}

	err = validation.Repeat(task.Repeat)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": err.Error()})
	}

	result := database.Gorm.Db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"Error": result.Error.Error()})
	} else if result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": models.ErrTaskNotFound})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{})
}
