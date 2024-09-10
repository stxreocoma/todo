package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stxreocoma/todo/database"
	"github.com/stxreocoma/todo/models"
	"github.com/stxreocoma/todo/utils"
)

func GetDate(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	now, err := time.Parse(utils.DateFormat, c.Query("now"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	date, err := utils.NextDate(now, c.Query("date"), c.Query("repeat"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString(date)
}

func GetTasks(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	tasks := make([]models.Task, 50)

	if c.Query("search") == "" {
		result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler").Limit(50).Find(&tasks)
		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"error": models.ErrRepeat.Error()})
		}
	} else {
		date, err := time.Parse("02.01.2006", c.Query("search"))
		if err != nil {
			result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? LIMIT 50", "%"+c.Query("search")+"%", "%"+c.Query("search")+"%").Find(&tasks)
			if result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"error": models.ErrSearch.Error()})
			}
		} else {
			result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT 50", date.Format(utils.DateFormat)).Find(&tasks)
			if result.Error != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"error": models.ErrSearch.Error()})
			}
		}
	}

	utils.Sort(tasks)

	return c.Status(fiber.StatusOK).JSON(map[string]any{"tasks": tasks})
}

func GetTask(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	var task models.Task

	if c.Query("id") == "" {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"error": models.ErrID.Error()})
	}

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"error": models.ErrTaskNotFound.Error()})
	}
	result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ? LIMIT 1", id).First(&task)
	if result.Error != nil || result.RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"error": models.ErrTaskNotFound.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(task)
}
