package main

import "strings"

func cleanInput(text string) []string {
	cleanedInput := []string{}
	tokens := strings.Split(text, " ")

	for i := 0; i < len(tokens); i++ {
		tokens[i] = strings.ToLower(tokens[i])
		tokens[i] = strings.TrimSpace(tokens[i])
		if len(tokens[i]) > 0 {
			cleanedInput = append(cleanedInput, tokens[i])
		}
	}

	return cleanedInput
}
