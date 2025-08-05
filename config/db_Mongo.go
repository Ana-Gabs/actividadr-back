package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

// InitMongoDB establece la conexión con MongoDB Atlas
func InitMongoDB() error {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env (se asumirá entorno de producción)")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return fmt.Errorf("falta la variable de entorno MONGODB_URI")
	}

	// Contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Opciones de conexión
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return fmt.Errorf("error conectando a MongoDB: %v", err)
	}

	// Verificar la conexión
	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("no se pudo hacer ping a MongoDB: %v", err)
	}

	MongoClient = client
	MongoDB = client.Database("actividadr")

	log.Println("Conexión exitosa con MongoDB Atlas")
	return nil
}

// CloseMongo cierra la conexión con MongoDB
func CloseMongo() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Fatalf("Error al cerrar conexión con MongoDB: %v", err)
		}
		log.Println("Conexión con MongoDB cerrada")
	}
}

// GetCollection devuelve una colección específica de la base de datos
func GetCollection(name string) *mongo.Collection {
	return MongoDB.Collection(name)
}
