package userRoute

import (
	userHandler "fleet/api/handlers/user"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoute(router fiber.Router) {
	user := router.Group("/user")
	user.Get("/getusers", userHandler.GetUsers)
	user.Get("/:userId", userHandler.GetUser)
	/*
	   user.Put("/:userId", userHandler.UpdateUser)
	   user.Delete("/:userId", userHandler.DeleteUser)
	*/
}
