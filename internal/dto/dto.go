package dto

type ChirpDto struct {
	Body string `json:"body"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ExpiresAt int    `json:"expires_in_seconds"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type ErrorDto struct {
	Error string `json:"error"`
}
