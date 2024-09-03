package auth

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var password = os.Getenv("TODO_PASSWORD")

func Authentication(next fiber.Handler) fiber.Handler {
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
				return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{"Error": "Authentication required"})
			}
		}
		return next(c)
	}
}

func Registration(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	var enteredPassword map[string]string

	err := c.BodyParser(&enteredPassword)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"Error": err.Error()})
	}

	if enteredPassword["password"] != password {
		return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{"Error": "Неверный пароль"})
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	token, err := jwtToken.SignedString([]byte(password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"Error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{"Token": token})
}
