package model

type RouteLogs struct {
	ID        int    `json:"id"`
	Lat       string `json:"lat"`
	Lng       string `json:"lng"`
	DateTime  string `json:"datetime"`
	VehicleID string `json:"vehicle_id"`
}
