package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/bogdan-cu/pokedexcli/internal/pokeapi"
)

func main() {
	config := pokeapi.Config{PrevUrl: "", NextUrl: locationAreaUrl}
	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		command := ""
		for key := range commands {
			if key == cleanedInput[0] {
				command = key
			}
		}
		if command == "" {
			fmt.Println("Unknown command")
		}

		if err := commands[command].callback(&config); err != nil {
			fmt.Printf("command execution failed: %s\n", err)
		}
	}
}
