// ./routes/user_routes.go

package routes

import (
	"github.com/Ana-Gabs/actividadr-back/controllers"
	"github.com/Ana-Gabs/actividadr-back/middlewares"
	"github.com/gofiber/fiber/v2"
)


func SetupUserRoutes(app *fiber.App) {
	// Rutas de autenticación
	app.Post("/login", middlewares.RateLimitMiddleware(), controllers.Login)
	app.Post("/register", middlewares.RateLimitMiddleware(), controllers.Register)
	app.Get("/info", middlewares.RateLimitMiddleware(), controllers.GetInfo)
	// app.Get("/info", middlewares.AuthMiddleware(), controllers.GetInfo) // Comentado como en el original
	//app.Post("/verify-otp", middlewares.RateLimitMiddleware(), controllers.VerifyOtp)
	app.Post("/verify-otp", controllers.VerifyOtp)
}
