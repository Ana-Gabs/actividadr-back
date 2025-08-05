// ./routes/logsss_routes.go
package routes

import (
	"github.com/Ana-Gabs/actividadr-back/controllers"
	"github.com/gofiber/fiber/v2"
	//"github.com/Ana-Gabs/actividadr-back/middlewares"
)

func SetupLogsRoutes(app *fiber.App) {
	// Rutas para obtener logs
	app.Get("/logs/level", controllers.GetLogsByLevel)
	app.Get("/logs/time", controllers.GetLogsByResponseTime)
	app.Get("/logs/status", controllers.GetLogsByStatus)

	// Rutas con rate limiting (descomentar para habilitar)
	// app.Get("/logs/level", middlewares.RateLimitMiddleware(), controllers.GetLogsByLevel)
	// app.Get("/logs/time", middlewares.RateLimitMiddleware(), controllers.GetLogsByResponseTime)
	// app.Get("/logs/status", middlewares.RateLimitMiddleware(), controllers.GetLogsByStatus)
}
