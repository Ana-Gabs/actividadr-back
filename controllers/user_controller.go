// ./controllers/user_controller.go

package controllers

import (
	"context"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Ana-Gabs/actividadr-back/config"
	"github.com/Ana-Gabs/actividadr-back/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func GetInfo(c *fiber.Ctx) error {

	if rand.Float64() < 0.3 {
		utils.LogAction("anonymous", "getInfo-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error interno del servidor"})
	}

	info := fiber.Map{
		"go_version": strings.TrimPrefix(runtime.Version(), "go"),
		"student": fiber.Map{
			"name":  "Ana Gabriela Contreras Jiménez",
			"group": "Grupo IDGS11",
		},
	}

	utils.LogAction("anonymous", "getInfo", "info")(c)
	return c.JSON(info)
}


func Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogAction("anonymous", "register-error", "error")(c)
		return c.Status(400).JSON(fiber.Map{"error": "Solicitud inválida"})
	}

	if req.Email == "" || req.Username == "" || req.Password == "" {
		utils.LogAction("anonymous", "register-error", "error")(c)
		return c.Status(400).JSON(fiber.Map{"error": "Todos los campos son obligatorios"})
	}

	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		utils.LogAction("anonymous", "register-error", "error")(c)
		return c.Status(400).JSON(fiber.Map{"error": "Email inválido"})
	}

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser bson.M
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err != mongo.ErrNoDocuments {
		utils.LogAction("anonymous", "register-error", "error")(c)
		return c.Status(400).JSON(fiber.Map{"error": "El usuario ya existe"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogAction("anonymous", "register-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error en el registro"})
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Actividadr-Back",
		AccountName: req.Email,
	})
	if err != nil {
		utils.LogAction("anonymous", "register-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error en el registro"})
	}

	_, err = collection.InsertOne(ctx, bson.M{
		"email":         req.Email,
		"username":      req.Username,
		"password":      string(hashedPassword),
		"mfa_secret":    key.Secret(),
		"mfaEnabled":    true,
		"date_register": time.Now(),
		"last_login":    nil,
	})
	if err != nil {
		utils.LogAction("anonymous", "register-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error en el registro"})
	}

	utils.LogAction(req.Email, "register", "info")(c)
	return c.Status(201).JSON(fiber.Map{
		"message":    "Usuario registrado con éxito",
		"mfa_secret": key.URL(),
	})
}

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		EmailOrUsername string `json:"emailOrUsername"`
		Password        string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogAction("anonymous", "login-error", "error")(c)
		return c.Status(400).JSON(fiber.Map{"error": "Solicitud inválida"})
	}

	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user bson.M
	err := collection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"email": req.EmailOrUsername},
			{"username": req.EmailOrUsername},
		},
	}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		utils.LogAction("anonymous", "login-error", "error")(c)
		return c.Status(401).JSON(fiber.Map{"error": "Credenciales incorrectas"})
	} else if err != nil {
		utils.LogAction("anonymous", "login-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error en el login"})
	}

	
	if err := bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(req.Password)); err != nil {
		utils.LogAction("anonymous", "login-error", "error")(c)
		return c.Status(401).JSON(fiber.Map{"error": "Credenciales incorrectas"})
	}

	
	if user["mfaEnabled"].(bool) {
		utils.LogAction(user["email"].(string), "login-mfa-required", "info")(c)
		return c.JSON(fiber.Map{
			"requiresMFA": true,
			"email":       user["email"],
		})
	}

	
	token, err := generateJWT(user["email"].(string))
	if err != nil {
		utils.LogAction("anonymous", "login-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error en el login"})
	}

	
	_, err = collection.UpdateOne(ctx, bson.M{"email": user["email"]}, bson.M{
		"$set": bson.M{"last_login": time.Now()},
	})
	if err != nil {
		utils.LogAction("anonymous", "login-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"error": "Error en el login"})
	}

	utils.LogAction(user["email"].(string), "login", "info")(c)
	return c.JSON(fiber.Map{"token": token})
}


func VerifyOtp(c *fiber.Ctx) error {
	type OtpRequest struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}

	var req OtpRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogAction("anonymous", "verifyOtp-error", "error")(c)
		return c.Status(400).JSON(fiber.Map{"message": "Faltan datos en la solicitud"})
	}

	
	collection := config.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user bson.M
	err := collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		utils.LogAction("anonymous", "verifyOtp-error", "error")(c)
		return c.Status(401).JSON(fiber.Map{"message": "Usuario no encontrado"})
	} else if err != nil {
		utils.LogAction("anonymous", "verifyOtp-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"message": "Error interno del servidor"})
	}

	if user["mfa_secret"] == nil {
		utils.LogAction("anonymous", "verifyOtp-error", "error")(c)
		return c.Status(400).JSON(fiber.Map{"message": "El usuario no tiene 2FA habilitado"})
	}

	
	isValid := totp.Validate(req.Token, user["mfa_secret"].(string))
	if !isValid {
		utils.LogAction(req.Email, "verifyOtp-error", "error")(c)
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Código OTP inválido o expirado",
		})
	}

	
	token, err := generateJWT(user["email"].(string))
	if err != nil {
		utils.LogAction("anonymous", "verifyOtp-error", "error")(c)
		return c.Status(500).JSON(fiber.Map{"message": "Error interno del servidor"})
	}

	utils.LogAction(req.Email, "verifyOtp-success", "info")(c)
	return c.JSON(fiber.Map{
		"success": true,
		"token":   token,
	})
}

func generateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
