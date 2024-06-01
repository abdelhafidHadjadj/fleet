package vehicleRouter

import (
	vehicleHandler "fleet/api/handlers/vehicle"

	"github.com/gofiber/fiber/v2"
)

func SetupVehicleRoute(router fiber.Router) {
	vehicle := router.Group("/vehicle")
	vehicle.Get("/getall", vehicleHandler.GetVehicles)
	vehicle.Get("/get/:vehicleId", vehicleHandler.GetVehicleByID)
	vehicle.Post("/add", vehicleHandler.CreateVehicle)
	vehicle.Patch("/update/:vehicleId", vehicleHandler.UpdateVehicle)
	vehicle.Delete("/delete/:vehicleId", vehicleHandler.DeleteVehicle)
	vehicle.Get("/drones", vehicleHandler.GetDrones)
	vehicle.Get("/drones_number", vehicleHandler.GetDronesNumber)
	vehicle.Get("/vehicles", vehicleHandler.GetVehicleWithoutDrones)
	vehicle.Get("/vehicles_number", vehicleHandler.GetVehicleNumber)
	vehicle.Get("/vehicles_status/:status", vehicleHandler.GetVehiclesNumberByStatus)
	vehicle.Get("/vehicles/created_by_month", vehicleHandler.GetVehiclesCreatedByMonth)

}
