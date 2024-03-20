package userRoute

import (
	userHandler "fleet/api/handlers/user"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoute(router fiber.Router) {
	user := router.Group("/user")
	user.Get("/getall", userHandler.GetUsers)
	user.Get("/get/:userId", userHandler.GetUserByID)
	user.Post("/add", userHandler.CreateUser)
	user.Patch("/update/:userId", userHandler.UpdateUser)
	user.Delete("/delete/:userId", userHandler.DeleteUser)
}
