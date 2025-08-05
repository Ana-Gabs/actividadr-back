package utils

import (
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/Ana-Gabs/actividadr-back/config"
)

// LogAction registra una acción HTTP en la colección "logs"
func LogAction(email string, action string, logLevel string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		duration := time.Since(start)

		hostname, _ := os.Hostname()

		status := c.Response().StatusCode()
		if logLevel == "" {
			if status >= 400 {
				logLevel = "error"
			} else {
				logLevel = "info"
			}
		}

		logEntry := map[string]interface{}{
			"email":         email,
			"action":        action,
			"logLevel":      logLevel,
			"timestamp":     time.Now(),
			"ip":            c.IP(),
			"userAgent":     c.Get("User-Agent", "Unknown"),
			"referer":       c.Get("Referer", "Unknown"),
			"origin":        c.Get("Origin", "Unknown"),
			"method":        c.Method(),
			"url":           c.OriginalURL(),
			"status":        status,
			"responseTime":  duration.Milliseconds(),
			"protocol":      c.Protocol(),
			"hostname":      hostname,
			"environment":   os.Getenv("NODE_ENV"),
			"goVersion":     strings.TrimPrefix(runtime.Version(), "go"),
			"pid":           os.Getpid(),
		}

		collection := config.GetCollection("logs")
		_, insertErr := collection.InsertOne(c.Context(), logEntry)
		if insertErr != nil {
			log.Println("Error al registrar log:", insertErr)
		}

		return err
	}
}
