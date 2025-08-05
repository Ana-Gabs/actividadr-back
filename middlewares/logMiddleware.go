// ./middleware/logMiddleware.go
package middlewares

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LogMiddleware registra información de cada solicitud HTTP
func LogMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Procesa la solicitud
	err := c.Next()

	// Calcula el tiempo de respuesta
	duration := time.Since(start)

	// Nivel de log
	status := c.Response().StatusCode()
	logLevel := "info"
	if status >= 400 {
		logLevel = "error"
	}

	// Datos del log
	logData := map[string]interface{}{
		"timestamp":     time.Now().Format(time.RFC3339),
		"method":        c.Method(),
		"url":           c.OriginalURL(),
		"status":        status,
		"response_time": duration.Milliseconds(),
		"ip":            c.IP(),
		"user_agent":    c.Get("User-Agent"),
	}

	log.Printf("[%s] %s %s %d %dms IP: %s UA: %s",
		logLevel,
		logData["method"],
		logData["url"],
		logData["status"],
		logData["response_time"],
		logData["ip"],
		logData["user_agent"],
	)

	return err
}
