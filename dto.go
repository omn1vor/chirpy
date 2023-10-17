package main

type chirpDto struct {
	Body string `json:"body"`
}

type userDto struct {
	Email string `json:"email"`
}

type errorDto struct {
	Error string `json:"error"`
}
