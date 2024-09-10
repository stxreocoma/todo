package auth

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var password = os.Getenv("TODO_PASSWORD")

func Authentication(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
		if len(password) > 0 {
			token := c.Cookies("token")
			if len(token) == 0 {
				return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{"error": "Authentication required"})
			}
			jwtToken := jwt.New(jwt.SigningMethodHS256)
			passwordToken, err := jwtToken.SignedString([]byte(password))
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{"error": "Authentication required"})
			} else if passwordToken != token {
				c.Status(fiber.StatusUnauthorized).JSON(map[string]any{"error": "Authentication required"})
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
		return c.Status(fiber.StatusBadRequest).JSON(map[string]any{"error": err.Error()})
	}

	if enteredPassword["password"] != password {

		return c.Status(fiber.StatusUnauthorized).JSON(map[string]any{"error": "Неверный пароль"})
	}

	jwtToken := jwt.New(jwt.SigningMethodHS256)
	token, err := jwtToken.SignedString([]byte(password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]any{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{"token": token})
}
