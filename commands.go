package main

import (
	"fmt"
	"github.com/bogdan-cu/pokedexcli/internal/pokecache"
	"os"

	"github.com/bogdan-cu/pokedexcli/internal/pokeapi"
)

const locationAreaUrl = "https://pokeapi.co/api/v2/location-area/"

type App struct {
	config *pokeapi.Config
	cache  *pokecache.Cache
}

type CliCommand struct {
	name        string
	description string
	callback    func(*App, ...string) error
}

var commands = map[string]CliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the CLI",
		callback:    (*App).commandExit,
	},
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    (*App).commandHelp,
	},
	"map": {
		name:        "map",
		description: "It displays the names of 20 location areas in the Pokemon world",
		callback:    (*App).commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "It displays the names of the previous 20 location areas in the Pokemon world",
		callback:    (*App).commandMapb,
	},
	"explore": {
		name:        "explore",
		description: "It returns a list of pokemon found in a given location",
		callback:    (*App).explore,
	},
}

func (a *App) commandExit(_ ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (a *App) commandHelp(_ ...string) error {
	fmt.Println(`
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: It displays the names of 20 location areas in the Pokemon world
mapb: It displays the names of the previous 20 location areas in the Pokemon world
explore <map_area>: It returns a list of pokemon found in a given location`)
	return nil
}

func (a *App) commandMap(_ ...string) error {
	url := a.config.NextUrl
	if cacheEntry, ok := a.cache.Get(url); ok {
		results := byteSliceToStringSlice(cacheEntry)
		_ = writeStrings(os.Stdout, results...)
		return nil
	}
	results, err := pokeapi.GetLocationArea(a.config, true)
	if err != nil {
		return err
	}
	_ = writeStrings(os.Stdout, results...)
	cacheEntry := stringSliceToByteSlice(results)
	a.cache.Add(url, cacheEntry)
	return nil
}

func (a *App) commandMapb(_ ...string) error {
	url := a.config.PrevUrl
	if cacheEntry, ok := a.cache.Get(url); ok {
		results := byteSliceToStringSlice(cacheEntry)
		_ = writeStrings(os.Stdout, results...)
		return nil
	}
	results, err := pokeapi.GetLocationArea(a.config, false)
	if err != nil {
		return err
	}
	_ = writeStrings(os.Stdout, results...)
	cacheEntry := stringSliceToByteSlice(results)
	a.cache.Add(url, cacheEntry)
	return nil
}

func (a *App) explore(locations ...string) error {
	if len(locations) == 0 {
		return fmt.Errorf("run the command on a location\n")
	}
	if len(locations) > 1 {
		return fmt.Errorf("one location at a time, please\n")
	}
	locationUrl := locationAreaUrl + locations[0]
	if results, ok := a.cache.Get(locationUrl); ok {
		pokemon := byteSliceToStringSlice(results)
		_ = writeStrings(os.Stdout, "Pokemon found in the area:")
		_ = writeStrings(os.Stdout, pokemon...)
	}
	pokemon, err := pokeapi.GetLocalPokemon(locationUrl)
	if err != nil {
		return err
	}
	a.cache.Add(locationUrl, stringSliceToByteSlice(pokemon))
	_ = writeStrings(os.Stdout, "Pokemon found in the area:")
	_ = writeStrings(os.Stdout, pokemon...)
	return nil
}
