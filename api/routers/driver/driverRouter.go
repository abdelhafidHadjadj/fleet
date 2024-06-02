package driverRouter

import (
	driverHandler "fleet/api/handlers/driver"

	"github.com/gofiber/fiber/v2"
)

func SetupDriverRoute(router fiber.Router) {
	driver := router.Group("/driver")
	driver.Get("/getall", driverHandler.GetDrivers)
	driver.Get("/get/:driverId", driverHandler.GetDriverByID)
	driver.Post("/add", driverHandler.CreateDriver)
	driver.Patch("/update/:driverId", driverHandler.UpdateDriver)
	driver.Delete("/delete/:driverId", driverHandler.DeleteDriver)
	driver.Get("/driver_number", driverHandler.GetDriversNumber)
}
