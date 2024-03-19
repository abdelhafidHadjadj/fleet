package model

type Route struct {
	ID            int    `json:"id"`
	Status        string `json:"status"`
	Departue_date string `json:"departue_date"`
	Arrival_date  string `json:"arrival_date"`
	DriverID      int    `json:"driver_id"`
	VehicleID     int    `json:"vehicle_id"`
	CreatedAt     string `json:"created_at"`
}
