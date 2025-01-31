package main

import (
	"bufio"
	"fmt"
	"github.com/bogdan-cu/pokedexcli/internal/pokecache"
	"os"
	"time"

	"github.com/bogdan-cu/pokedexcli/internal/pokeapi"
)

const locationAreaUrl = "https://pokeapi.co/api/v2/location-area/"

func main() {
	app := App{
		config: &pokeapi.Config{PrevUrl: "", NextUrl: locationAreaUrl},
		cache:  pokecache.NewCache(5 * time.Second),
	}
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

		if err := commands[command].callback(&app); err != nil {
			fmt.Printf("command execution failed: %s\n", err)
		}
	}
}
