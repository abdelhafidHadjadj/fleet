package model

type Route struct {
	ID            int    `json:"id"`
	Status        string `json:"status"`
	Departue_date string `json:"departure_date"`
	Arrival_date  string `json:"arrival_date"`
	DriverID      string `json:"driver_id"`
	VehicleID     string `json:"vehicle_id"`
	CreatedAt     string `json:"created_at"`
	CreatedBy     string `json:"created_by"`
}
