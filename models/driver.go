package model

type Driver struct {
	ID              int    `json:"id"`
	Register_number string `json:"register_number"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Date_of_birth   string `json:"date_of_birth"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	Class           string `json:"class"`
	Status          string `json:"status"`
	CreatedBy       int    `json:"created_by"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
