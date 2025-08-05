// ./controllers/logs_controller.go

package controllers

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/Ana-Gabs/actividadr-back/config"
	"github.com/Ana-Gabs/actividadr-back/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)


func GetLogsByLevel(c *fiber.Ctx) error {
	collection := config.GetCollection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		utils.LogAction("anonymous", "getLogsByLevel-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener los logs por nivel"})
	}
	defer cursor.Close(ctx)

	groupedByLevel := make(map[string]int)
	for cursor.Next(ctx) {
		var log bson.M
		if err := cursor.Decode(&log); err != nil {
			utils.LogAction("anonymous", "getLogsByLevel-error", "error")(c)
			return c.Status(500).JSON(fiber.Map{"error": "Error al procesar los logs"})
		}
		level, ok := log["logLevel"].(string)
		if !ok {
			level = "unknown"
		}
		groupedByLevel[level]++
	}

	if err := cursor.Err(); err != nil {
		utils.LogAction("anonymous", "getLogsByLevel-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener los logs por nivel"})
	}

	// Registrar acción
	utils.LogAction("anonymous", "getLogsByLevel", "info")(c)
	return c.Status(200).JSON(groupedByLevel)
}

// GetLogsByResponseTime agrupa los logs por tiempo de respuesta
func GetLogsByResponseTime(c *fiber.Ctx) error {
	collection := config.GetCollection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Obtener todos los logs
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		utils.LogAction("anonymous", "getLogsByResponseTime-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener los logs por tiempo de respuesta"})
	}
	defer cursor.Close(ctx)

	// Agrupar logs por tiempo de respuesta (en segundos)
	responseTimeStats := make(map[int]int)
	for cursor.Next(ctx) {
		var log bson.M
		if err := cursor.Decode(&log); err != nil {
			utils.LogAction("anonymous", "getLogsByResponseTime-error", "error")(c)
			return c.Status(500).JSON(fiber.Map{"error": "Error al procesar los logs"})
		}
		responseTime, ok := log["responseTime"].(float64)
		if !ok {
			responseTime = 0
		}
		rangeKey := int(math.Floor(responseTime))
		responseTimeStats[rangeKey]++
	}

	if err := cursor.Err(); err != nil {
		utils.LogAction("anonymous", "getLogsByResponseTime-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener los logs por tiempo de respuesta"})
	}

	// Registrar acción
	utils.LogAction("anonymous", "getLogsByResponseTime", "info")(c)
	return c.Status(200).JSON(responseTimeStats)
}

// GetLogsByStatus agrupa los logs por código de estado HTTP
func GetLogsByStatus(c *fiber.Ctx) error {
	collection := config.GetCollection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Obtener todos los logs
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		utils.LogAction("anonymous", "getLogsByStatus-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener los logs por código de estado"})
	}
	defer cursor.Close(ctx)

	// Agrupar logs por status
	groupedByStatus := make(map[string]int)
	for cursor.Next(ctx) {
		var log bson.M
		if err := cursor.Decode(&log); err != nil {
			utils.LogAction("anonymous", "getLogsByStatus-error", "error")(c)
			return c.Status(500).JSON(fiber.Map{"error": "Error al procesar los logs"})
		}
		var status string
		if log["status"] != nil {
			// Convertir status a string, ya que puede ser int32 en MongoDB
			switch v := log["status"].(type) {
			case int32:
				status = fmt.Sprintf("%d", v)
			case int64:
				status = fmt.Sprintf("%d", v)
			case float64:
				status = fmt.Sprintf("%d", int(v))
			default:
				status = "unknown"
			}
		} else {
			status = "unknown"
		}
		groupedByStatus[status]++
	}

	if err := cursor.Err(); err != nil {
		utils.LogAction("anonymous", "getLogsByStatus-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener los logs por código de estado"})
	}

	// Registrar acción
	utils.LogAction("anonymous", "getLogsByStatus", "info")(c)
	return c.Status(200).JSON(groupedByStatus)
}
