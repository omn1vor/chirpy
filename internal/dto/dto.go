package dto

type ChirpDto struct {
	Body string `json:"body"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type ErrorDto struct {
	Error string `json:"error"`
}
