package model

type User struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     int    `json:"phone"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}
