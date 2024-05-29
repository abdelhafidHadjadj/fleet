package model

type LoginRequest struct {
	LoginType string `json:"loginType"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
