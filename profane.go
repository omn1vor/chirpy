package main

import (
	"strings"
)

func profaneWords() map[string]bool {
	words := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	return words
}

func replaceProfanity(body string) string {
	const profaneMask = "****"
	badWords := profaneWords()
	allWords := strings.Split(body, " ")
	for i, word := range allWords {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			allWords[i] = profaneMask
		}
	}
	return strings.Join(allWords, " ")
}
