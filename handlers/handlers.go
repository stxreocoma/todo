package handlers

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stxreocoma/todo/database"
	"github.com/stxreocoma/todo/models"
	"github.com/stxreocoma/todo/utils"
	"github.com/stxreocoma/todo/validation"
)

var password = os.Getenv("TODO_PASSWORD")

func Date(c *fiber.Ctx) error {
	now, err := time.Parse(utils.DateFormat, c.Query("now"))
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	date, err := utils.NextDate(now, c.Query("date"), c.Query("repeat"))
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString(date)
}

func PostTask(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	date, err := time.Parse(utils.DateFormat, task.Date)
	if err != nil {
		if task.Date == "" {
			date = time.Now()
		} else {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
		}
	} else if date.Unix() < time.Now().Unix() && (time.Now().Unix()-date.Unix() >= 86400 || time.Now().Day() != date.Day()) {
		if task.Repeat == "" {
			date = time.Now()
		} else {
			dateString, err := utils.NextDate(time.Now(), date.Format(utils.DateFormat), task.Repeat)
			if err != nil {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
			}
			date, err = time.Parse(utils.DateFormat, dateString)
			if err != nil {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
			}
		}
	}
	task.Date = date.Format(utils.DateFormat)

	if len(task.Title) == 0 {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "no title"})
	}

	err = validation.Repeat(task.Repeat)
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	result := database.Gorm.Db.Omit("created_at", "updated_at", "deleted_at").Create(&task)
	if result.Error != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: result.Error.Error()})
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(fiber.StatusOK).JSON(models.CorrectResponse{ID: strconv.Itoa(task.ID)})
}

func GetTasks(c *fiber.Ctx) error {
	tasks := make([]models.TaskForTests, 50)

	if c.Query("search") == "" {
		result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler LIMIT 50").Limit(50).Find(&tasks)
		if result.Error != nil {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
		}
	} else {
		date, err := time.Parse("02.01.2006", c.Query("search"))
		if err != nil {
			result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? LIMIT 50", "%"+c.Query("search")+"%").Find(&tasks)
			if result.Error != nil {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "wrong search format"})
			}
		} else {
			result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT 50", date.Format(utils.DateFormat)).Find(&tasks)
			if result.Error != nil {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "wrong search format"})
			}
		}
	}

	for i := 0; i < len(tasks)-1; i++ {
		for j := i + 1; j < len(tasks); j++ {
			aDate, _ := time.Parse(utils.DateFormat, tasks[i].Date)
			bDate, _ := time.Parse(utils.DateFormat, tasks[j].Date)
			if aDate.Unix() > bDate.Unix() {
				tasks[i], tasks[j] = tasks[j], tasks[i]
			}
		}
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(fiber.StatusOK).JSON(models.CorrectGetTasksForTests{Tasks: tasks})
}

func GetTask(c *fiber.Ctx) error {
	var task models.Task

	if c.Query("id") == "" {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Не указан идентификатор"})
	}

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Задача не найдена"})
	}
	result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ? LIMIT 1", id).First(&task)
	if result.Error != nil || result.RowsAffected == 0 {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Задача не найдена"})
	}

	taskForTests := models.TaskForTests{
		ID:      c.Query("id"),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(fiber.StatusOK).JSON(taskForTests)
}

func UpdateTask(c *fiber.Ctx) error {
	var taskForTests models.TaskForTests
	var task models.Task

	if err := c.BodyParser(&taskForTests); err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Задача не найдена"})
	}

	id, err := strconv.Atoi(taskForTests.ID)
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Задача не найдена"})
	}

	task.ID = id
	task.Date = taskForTests.Date
	task.Title = taskForTests.Title
	task.Comment = taskForTests.Comment
	task.Repeat = taskForTests.Repeat

	date, err := time.Parse(utils.DateFormat, task.Date)
	if err != nil {
		if task.Date == "" {
			date = time.Now()
		} else {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
		}
	} else if date.Unix() < time.Now().Unix() {
		if task.Repeat == "" {
			date = time.Now()
		} else {
			dateString, err := utils.NextDate(time.Now(), date.Format(utils.DateFormat), task.Repeat)
			if err != nil {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
			}
			date, err = time.Parse(utils.DateFormat, dateString)
			if err != nil {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
			}
		}
	}
	task.Date = date.Format(utils.DateFormat)

	if len(task.Title) == 0 {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "no title"})
	}

	err = validation.Repeat(task.Repeat)
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	result := database.Gorm.Db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if result.Error != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: result.Error.Error()})
	} else if result.RowsAffected == 0 {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Задача не найдена"})
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(fiber.StatusOK).JSON(models.CorrectResponse{})
}

func DoneTask(c *fiber.Ctx) error {
	var task models.Task

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ? LIMIT 1", id).First(&task)
	if result.Error != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: result.Error.Error()})
	}

	if task.Repeat == "" {
		result = database.Gorm.Db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if result.Error != nil {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: result.Error.Error()})
		}

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusOK).JSON(models.CorrectResponse{})
	} else {
		newDate, err := utils.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
		}

		result = database.Gorm.Db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, id)
		if result.Error != nil {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: result.Error.Error()})
		}
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(fiber.StatusOK).JSON(models.CorrectResponse{})
}

func DeleteTask(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	result := database.Gorm.Db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if result.Error != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: result.Error.Error()})
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(fiber.StatusOK).JSON(models.CorrectResponse{})
}

func Auth(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(password) > 0 {
			var token string

			if len(c.Cookies("token")) != 0 {
				token = c.Cookies("token")
				log.Println("No token")
			}
			var valid bool
			jwtToken := jwt.New(jwt.SigningMethodHS256)
			passwordToken, err := jwtToken.SignedString([]byte(password))
			if err != nil {
				valid = false
			} else if passwordToken == token {
				valid = true
			}

			if !valid {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Authentification required"})
			}
		}
		return next(c)
	}
}

func Registration(c *fiber.Ctx) error {
	var enteredPassword map[string]string

	err := c.BodyParser(&enteredPassword)
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	if enteredPassword["password"] != password {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Неверный пароль"})
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	token, err := jwtToken.SignedString([]byte(password))
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(fiber.StatusOK).JSON(models.CorrectAuth{Token: token})
}
