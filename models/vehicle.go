package model

type Vehicle struct {
	ID               int    `json:"id"`
	Register_number  string `json:"register_number"`
	Name             string `json:"name"`
	Model            string `json:"model"`
	Type             string `json:"type"`
	Type_charge      string `json:"type_charge"`
	Currrent_charge  string `json:"currrent_charge"`
	Charge_capacity  string `json:"charge_capacity"`
	Current_distance string `json:"current_distance"`
	Current_position string `json:"current_position"`
	Status           string `json:"status"`
	Connection_key   string `json:"connection_key"`
	CreatedAt        string `json:"created_At"`
	CreatedBy        string `json:"created_by"`
}
