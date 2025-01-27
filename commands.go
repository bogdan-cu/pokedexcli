package main

import (
	"fmt"
	"os"

	"github.com/bogdan-cu/pokedexcli/internal/pokeapi"
)

const locationAreaUrl = "https://pokeapi.co/api/v2/location-area/"

type CliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config) error
}

var commands = map[string]CliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the CLI",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "It displays the names of 20 location areas in the Pokemon world",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "It displays the names of the previous 20 location areas in the Pokemon world",
		callback:    commandMapb,
	},
}

func commandExit(config *pokeapi.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *pokeapi.Config) error {
	fmt.Println(`
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: It displays the names of 20 location areas in the Pokemon world
mapb: It displays the names of the previous 20 location areas in the Pokemon world`)
	return nil
}

func commandMap(config *pokeapi.Config) error {
	results, err := pokeapi.GetLocationArea(config, true)
	if err != nil {
		return err
	}
	for _, result := range results {
		fmt.Println(result)
	}
	return nil
}

func commandMapb(config *pokeapi.Config) error {
	results, err := pokeapi.GetLocationArea(config, false)
	if err != nil {
		return err
	}
	for _, result := range results {
		fmt.Println(result)
	}
	return nil
}
