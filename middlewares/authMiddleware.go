// ./middleware/authMiddleware.go
package middlewares

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)


func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Token inválido",
		})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")


	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Token inválido o expirado",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		c.Locals("user", claims)
	}

	return c.Next()
}

