package authHandler

import (
	config "fleet/config"
	"fleet/database"
	models "fleet/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	// Extract the credentials from the request body
	loginRequest := new(models.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body: " + err.Error(),
		})
	}

	// Connect to the database
	db := database.ConnectionDB()

	// Variables to store user information
	var user models.User
	var driver models.Driver
	var hashedPassword string

	if loginRequest.LoginType == "user" {
		query := "SELECT id, role, firstname, lastname, email, password FROM USER WHERE email = ?"
		row := db.QueryRow(query, loginRequest.Email)
		if err := row.Scan(&user.ID, &user.Role, &user.Firstname, &user.Lastname, &user.Email, &user.Password); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
		hashedPassword = user.Password
	} else if loginRequest.LoginType == "driver" {
		query := "SELECT id, register_number, firstname, lastname, email, password FROM DRIVER WHERE email = ?"
		row := db.QueryRow(query, loginRequest.Email)
		if err := row.Scan(&driver.ID, &driver.Register_number, &driver.Firstname, &driver.Lastname, &driver.Email, &driver.Password); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
		hashedPassword = driver.Password
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid login type",
		})
	}

	// Check the hashed password
	if !CheckPasswordHash(loginRequest.Password, hashedPassword) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Create the JWT claims, which includes the user ID and expiry time
	day := time.Hour * 24
	claims := jwt.MapClaims{
		"exp": time.Now().Add(day * 1).Unix(),
	}

	if loginRequest.LoginType == "user" {
		claims["ID"] = user.ID
		claims["firstname"] = user.Firstname
		claims["lastname"] = user.Lastname
		claims["email"] = user.Email
		claims["role"] = user.Role
	} else if loginRequest.LoginType == "driver" {
		claims["ID"] = driver.ID
		claims["firstname"] = driver.Firstname
		claims["lastname"] = driver.Lastname
		claims["email"] = driver.Email
		claims["role"] = "driver"
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response
	t, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token: " + err.Error(),
		})
	}

	// Return the token
	return c.JSON(models.LoginResponse{
		Token: t,
	})
}

// Protected route
func Protected(c *fiber.Ctx) error {
	// Get the user from the context and return it
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	// Get common claims
	email := claims["email"].(string)
	role := claims["role"].(string)

	// Create response message based on role
	var responseMessage string
	if role == "operator" {
		firstname := claims["firstname"].(string)
		lastname := claims["lastname"].(string)
		responseMessage = "Welcome 👋 User " + firstname + " " + lastname + " (" + email + ")"
	} else if role == "driver" {
		firstname := claims["firstname"].(string)
		lastname := claims["lastname"].(string)
		responseMessage = "Welcome 👋 Driver " + firstname + " " + lastname
	} else if role == "admin" {
		firstname := claims["firstname"].(string)
		lastname := claims["lastname"].(string)
		responseMessage = "Welcome 👋 admin " + firstname + " " + lastname + " (" + email + ")"
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid role in token",
		})
	}

	return c.SendString(responseMessage)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
