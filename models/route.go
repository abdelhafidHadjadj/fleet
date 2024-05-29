package model

type Route struct {
	ID                    int    `json:"id"`
	Status                string `json:"status"`
	Departue_date         string `json:"departure_date"`
	Arrival_date          string `json:"arrival_date"`
	Lat_start             string `json:"lat_start"`
	Lng_start             string `json:"lng_start"`
	Lat_end               string `json:"lat_end"`
	Lng_end               string `json:"lng_end"`
	Departure_city        string `json:"departure_city"`
	Arrival_city          string `json:"arrival_city"`
	DriverID              string `json:"driver_id"`
	DriverName            string `json:"driver_name"`
	VehicleID             string `json:"vehicle_id"`
	VehicleRegisterNumber string `json:"veh_register_number"`
	VehicleName           string `json:"vehicle_name"`
	CreatedAt             string `json:"created_at"`
	CreatedBy             string `json:"created_by"`
}
