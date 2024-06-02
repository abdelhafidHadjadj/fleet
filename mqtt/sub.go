// sub.go
package mqtt

import (

	"fleet/database"

	"fmt"
	"log"

)

func logToDb(vehicleID, lat, lng string) {
	fmt.Printf("id: %s", vehicleID)
	fmt.Printf("lat: %s", lat)
	fmt.Printf("lng: %s", lng)
	db := database.ConnectionDB()
	stmt, err := db.Prepare("INSERT INTO ROUTE_LOGS (vehicle_id, lat, lng) VALUES (?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(vehicleID, lat, lng)
	if err != nil {
		log.Fatal(err)
	}
}
