package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func cleanInput(text string) []string {
	var cleanedInput []string
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

func writeStrings(w io.Writer, entries []string) error {
	for _, entry := range entries {
		_, err := fmt.Fprintf(w, entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func stringSliceToByteSlice(input []string) []byte {
	var b [][]byte
	for _, i := range input {
		b = append(b, []byte(i))
	}
	return bytes.Join(b, []byte("\n"))
}

func byteSliceToStringSlice(input []byte) []string {
	var s []string
	rawRows := bytes.Split(input, []byte("\n"))
	for _, row := range rawRows {
		s = append(s, string(row))
	}
	return s
}
