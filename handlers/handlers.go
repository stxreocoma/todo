package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stxreocoma/todo/database"
	"github.com/stxreocoma/todo/models"
)

var monthsMap = map[string]string{
	"1":  "January",
	"2":  "February",
	"3":  "March",
	"4":  "April",
	"5":  "May",
	"6":  "June",
	"7":  "July",
	"8":  "August",
	"9":  "September",
	"10": "October",
	"11": "November",
	"12": "December",
}

func Index(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./web/index.html"))
}

func lenMonth(date time.Time) int {
	switch {
	case date.Month().String() == monthsMap["1"] || date.Month().String() == monthsMap["3"] || date.Month().String() == monthsMap["5"] || date.Month().String() == monthsMap["7"] || date.Month().String() == monthsMap["8"] || date.Month().String() == monthsMap["10"] || date.Month().String() == monthsMap["12"]:
		return 31
	case date.Month().String() == monthsMap["2"]:
		if (date.Year()%4 == 0 && date.Year()%100 != 0) || date.Year()%400 == 0 {
			return 29
		} else {
			return 28
		}
	case date.Month().String() == monthsMap["4"] || date.Month().String() == monthsMap["6"] || date.Month().String() == monthsMap["9"] || date.Month().String() == monthsMap["11"]:
		return 30
	}
	return 0
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	d, err := time.Parse("20060102", date)
	if err != nil {
		log.Println(err)
		return "", err
	}

	params := strings.Split(repeat, " ")
	if len(params) == 1 && params[0] != "y" {
		return "", fmt.Errorf("wrong parameters")
	}

	switch params[0] {
	case "y":
		for d.Unix() <= now.Unix() {
			d = d.AddDate(1, 0, 0)
		}

		if d.Format("20060102") == date {
			d = d.AddDate(1, 0, 0)
		}

		return d.Format("20060102"), nil

	case "d":
		param, err := strconv.Atoi(params[1])
		if err != nil {
			return "", err
		} else if param > 400 {
			return "", fmt.Errorf("wrong number: %d\nmax number: 400", param)
		}
		for d.Unix() <= now.Unix() {
			d = d.AddDate(0, 0, param)
		}

		if d.Format("20060102") == date {
			d = d.AddDate(0, 0, param)
		}

		return d.Format("20060102"), nil

	case "w":
		days := strings.Split(params[1], ",")

		for {
			d = d.AddDate(0, 0, 1)
			if d.Unix() > now.Unix() {
				for _, v := range days {
					day, err := strconv.Atoi(v)
					if err != nil {
						return "", err
					} else if day > 7 {
						return "", fmt.Errorf("wrong number: %d\nmax number: 7", day)
					}
					weekDay := int(d.Weekday())
					if weekDay == 0 {
						weekDay += 7
					}
					log.Println((int(d.Weekday())), day)
					if weekDay == day {
						return d.Format("20060102"), nil
					}
				}
			}
		}

	case "m":
		days := strings.Split(params[1], ",")

		if len(params) == 2 {
			for {
				d = d.AddDate(0, 0, 1)
				if d.Unix() > now.Unix() {
					for _, v := range days {
						day, err := strconv.Atoi(v)
						if err != nil {
							return "", err
						} else if day > 31 || day < -2 {
							return "", fmt.Errorf("wrong number: %d\nmax number: 31\nmin number: -2", day)
						}

						if d.Day() == day || d.Day() == lenMonth(d)+day+1 {
							return d.Format("20060102"), nil
						}
					}
				}
			}
		} else {
			months := strings.Split(params[2], ",")

			for {
				d = d.AddDate(0, 0, 1)
				if d.Unix() > now.Unix() {
					for _, v1 := range months {
						month, err := strconv.Atoi(v1)
						if err != nil {
							return "", nil
						}

						if month > 12 || month < 1 {
							return "", fmt.Errorf("wrong number: %d\nmax number: 12", month)
						}

						for _, v2 := range days {
							day, err := strconv.Atoi(v2)
							if err != nil {
								return "", err
							} else if day > 31 || day < -2 {
								return "", fmt.Errorf("wrong number: %d\nmax number: 31\nmin number: -2", day)
							}

							if (d.Day() == day || d.Day() == lenMonth(d)+day+1) && int(d.Month()) == month {
								return d.Format("20060102"), nil
							}
						}
					}
				}
			}
		}

	default:
		return "", fmt.Errorf("wrong repeat format")
	}
}

func Date(c *fiber.Ctx) error {
	log.Println("params: ", c.Query("now"), c.Query("date"), c.Query("repeat"))
	now, err := time.Parse("20060102", c.Query("now"))
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	date, err := NextDate(now, c.Query("date"), c.Query("repeat"))
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	log.Println(date)
	return c.Status(fiber.StatusOK).SendString(date)
}

func PostTask(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	log.Println("params: ", task.Date, task.Title, task.Comment, task.Repeat)

	date, err := time.Parse("20060102", task.Date)
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
			dateString, err := NextDate(time.Now(), date.Format("20060102"), task.Repeat)
			if err != nil {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
			}
			date, err = time.Parse("20060102", dateString)
			if err != nil {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
			}
		}
	}
	task.Date = date.Format("20060102")

	if len(task.Title) == 0 {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "no title"})
	}

	if len(task.Repeat) > 0 && task.Repeat != "y" {
		repeat := strings.Split(task.Repeat, " ")
		switch repeat[0] {
		case "d":
			days := strings.Split(repeat[1], ",")
			for _, v := range days {
				day, err := strconv.Atoi(v)
				if err != nil {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				} else if day < 1 || day > 400 {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				}
			}
		case "w":
			days := strings.Split(repeat[1], ",")
			for _, v := range days {
				day, err := strconv.Atoi(v)
				if err != nil {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				} else if day < 1 || day > 7 {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				}
			}
		case "m":
			days := strings.Split(repeat[1], ",")
			for _, v := range days {
				day, err := strconv.Atoi(v)
				if err != nil {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				} else if day < -2 || day > 31 {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				}
			}
			if len(repeat) == 2 {
				months := strings.Split(repeat[2], ",")
				for _, month := range months {
					if _, ok := monthsMap[month]; !ok {
						c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
						return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
					}
				}
			}
		default:
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
		}
	}

	result := database.Gorm.Db.Omit("created_at", "updated_at", "deleted_at").Create(&task)
	if result.Error != nil {
		log.Println(err.Error())
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: result.Error.Error()})
	}

	log.Println(task.ID)
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(fiber.StatusOK).JSON(models.CorrectResponse{ID: strconv.Itoa(task.ID)})
}

func GetTasks(c *fiber.Ctx) error {
	tasks := make([]models.TaskForTests, 50)

	log.Println("search: ", c.Query("search"))

	if c.Query("search") == "" {
		result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler LIMIT 50").Limit(50).Find(&tasks)
		if result.Error != nil {
			log.Println(result.Error.Error())
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
		}
	} else {
		date, err := time.Parse("02.01.2006", c.Query("search"))
		if err != nil {
			result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? LIMIT 50", "%"+c.Query("search")+"%").Find(&tasks)
			log.Println("rows by search: ", result.RowsAffected)
			if result.Error != nil {
				log.Println(result.Error.Error())
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "wrong search format"})
			}
		} else {
			result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT 50", date.Format("20060102")).Find(&tasks)
			log.Println("rows by search: ", result.RowsAffected)
			if result.Error != nil {
				log.Println(result.Error.Error())
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "wrong search format"})
			}
		}
	}

	for i := 0; i < len(tasks)-1; i++ {
		for j := i + 1; j < len(tasks); j++ {
			aDate, _ := time.Parse("20060102", tasks[i].Date)
			bDate, _ := time.Parse("20060102", tasks[j].Date)
			if aDate.Unix() > bDate.Unix() {
				tasks[i], tasks[j] = tasks[j], tasks[i]
			}
		}
	}

	log.Println(tasks)
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
		log.Println(result.Error.Error())
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
		log.Println(err.Error())
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

	log.Println("params: ", task.Date, task.Title, task.Comment, task.Repeat)

	date, err := time.Parse("20060102", task.Date)
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
			dateString, err := NextDate(time.Now(), date.Format("20060102"), task.Repeat)
			if err != nil {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
			}
			date, err = time.Parse("20060102", dateString)
			if err != nil {
				c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
			}
		}
	}
	task.Date = date.Format("20060102")

	if len(task.Title) == 0 {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "no title"})
	}

	if len(task.Repeat) > 0 && task.Repeat != "y" {
		repeat := strings.Split(task.Repeat, " ")
		switch repeat[0] {
		case "d":
			days := strings.Split(repeat[1], ",")
			for _, v := range days {
				day, err := strconv.Atoi(v)
				if err != nil {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				} else if day < 1 || day > 400 {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				}
			}
		case "w":
			days := strings.Split(repeat[1], ",")
			for _, v := range days {
				day, err := strconv.Atoi(v)
				if err != nil {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				} else if day < 1 || day > 7 {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				}
			}
		case "m":
			days := strings.Split(repeat[1], ",")
			for _, v := range days {
				day, err := strconv.Atoi(v)
				if err != nil {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				} else if day < -2 || day > 31 {
					c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
					return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
				}
			}
			if len(repeat) == 2 {
				months := strings.Split(repeat[2], ",")
				for _, month := range months {
					if _, ok := monthsMap[month]; !ok {
						c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
						return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
					}
				}
			}
		default:
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "wrong repeat format"})
		}
	}

	result := database.Gorm.Db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if result.Error != nil {
		log.Println(err.Error())
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: result.Error.Error()})
	} else if result.RowsAffected == 0 {
		log.Println("Задача не найдена")
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "Задача не найдена"})
	}

	log.Println(task.ID)
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(fiber.StatusOK).JSON(models.CorrectResponse{})
}

func DoneTask(c *fiber.Ctx) error {
	var task models.Task

	log.Println(c.Query("id"))
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	result := database.Gorm.Db.Raw("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ? LIMIT 1", id).First(&task)
	if result.Error != nil {
		log.Println(result.Error.Error())
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: result.Error.Error()})
	}

	if task.Repeat == "" {
		result = database.Gorm.Db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if result.Error != nil {
			log.Println(result.Error.Error())
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: result.Error.Error()})
		}

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusOK).JSON(models.CorrectResponse{})
	} else {
		newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Println(err.Error())
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
		}

		result = database.Gorm.Db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, id)
		if result.Error != nil {
			log.Println(result.Error.Error())
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: result.Error.Error()})
		}
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(fiber.StatusOK).JSON(models.CorrectResponse{})
}

func DeleteTask(c *fiber.Ctx) error {
	log.Println(c.Query("id"))
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
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var token string

			if len(c.Cookies("token")) != 0 {
				token = c.Cookies("token")
				log.Println("No token")
			}
			var valid bool
			log.Println("token:", token)
			jwtToken := jwt.New(jwt.SigningMethodHS256)
			passwordToken, err := jwtToken.SignedString([]byte(pass))
			if err != nil {
				valid = false
				log.Println("validation:", err)
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
	password := os.Getenv("TODO_PASSWORD")

	var enteredPassword map[string]string

	err := c.BodyParser(&enteredPassword)
	log.Println("entered:", enteredPassword["password"], "real:", password)
	if err != nil {
		log.Println("parsing:", err)
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: err.Error()})
	}

	if enteredPassword["password"] != password {
		log.Println("мимо")
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "Неверный пароль"})
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	token, err := jwtToken.SignedString([]byte(password))
	if err != nil {
		log.Println("creating token:", err)
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	log.Println(token)
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Status(fiber.StatusOK).JSON(models.CorrectAuth{Token: token})
}
