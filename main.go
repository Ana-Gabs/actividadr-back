package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Ana-Gabs/actividadr-back/config"
	"github.com/Ana-Gabs/actividadr-back/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

// main inicializa y arranca el servidor
func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al cargar el archivo .env:", err)
	}

	// Inicializar conexión con MongoDB
	if err := config.InitMongoDB(); err != nil {
		log.Fatal("Error al inicializar MongoDB:", err)
	}
	defer config.CloseMongo() // Cerrar conexión al finalizar

	// Obtener variables de entorno
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Puerto por defecto si no está definido
	}
	ipWebserviceURL := os.Getenv("IP_WEBSERVICE_URL")
	if ipWebserviceURL == "" {
		ipWebserviceURL = "localhost" // Valor por defecto
	}

	// Inicializar la aplicación Fiber
	app := fiber.New()

	// Middlewares
	app.Use(logger.New()) // Reemplazo de logMiddleware
	app.Use(cors.New())   // Habilitar CORS

	// Configurar rutas
	routes.SetupUserRoutes(app)
	routes.SetupLogsRoutes(app)

	// Verificar la conexión con MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collections, err := config.MongoClient.Database("actividadr-back").ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Fatal("Error al conectar con MongoDB:", err)
	}
	fmt.Println("Conexión con MongoDB establecida correctamente. Colecciones encontradas:", collections)

	// Iniciar el servidor
	listenAddr := fmt.Sprintf("%s:%s", ipWebserviceURL, port)
	if err := app.Listen(listenAddr); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
	fmt.Printf("Servidor escuchando en http://%s\n", listenAddr)
}
