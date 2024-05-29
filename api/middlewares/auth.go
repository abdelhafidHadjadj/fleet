package middlewares

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(secret)},
		TokenLookup:  "header:Authorization",
		AuthScheme:   "Bearer",
		ContextKey:   "user",
		ErrorHandler: jwtError,
	})
}
func jwtError(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Next()
}
