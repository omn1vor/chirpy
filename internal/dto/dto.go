package dto

type ChirpDto struct {
	Body     string `json:"body"`
	AuthorId int    `json:"aouthor_id"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
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
