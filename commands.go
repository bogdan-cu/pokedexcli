package main

import (
	"fmt"
	"github.com/bogdan-cu/pokedexcli/internal/pokecache"
	"os"

	"github.com/bogdan-cu/pokedexcli/internal/pokeapi"
)

type App struct {
	config *pokeapi.Config
	cache  *pokecache.Cache
}

type CliCommand struct {
	name        string
	description string
	callback    func(*App) error
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
}

func (a *App) commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (a *App) commandHelp() error {
	fmt.Println(`
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: It displays the names of 20 location areas in the Pokemon world
mapb: It displays the names of the previous 20 location areas in the Pokemon world`)
	return nil
}

func (a *App) commandMap() error {
	url := a.config.NextUrl
	if cacheEntry, ok := a.cache.Get(url); ok {
		results := byteSliceToStringSlice(cacheEntry)
		err := writeStrings(os.Stdout, results)
		if err != nil {
			return err
		}
		return nil
	}
	results, err := pokeapi.GetLocationArea(a.config, true)
	if err != nil {
		return err
	}
	err = writeStrings(os.Stdout, results)
	if err != nil {
		return err
	}
	cacheEntry := stringSliceToByteSlice(results)
	a.cache.Add(url, cacheEntry)
	return nil
}

func (a *App) commandMapb() error {
	url := a.config.PrevUrl
	if cacheEntry, ok := a.cache.Get(url); ok {
		results := byteSliceToStringSlice(cacheEntry)
		err := writeStrings(os.Stdout, results)
		if err != nil {
			return err
		}
		return nil
	}
	results, err := pokeapi.GetLocationArea(a.config, false)
	if err != nil {
		return err
	}
	err = writeStrings(os.Stdout, results)
	if err != nil {
		return err
	}
	cacheEntry := stringSliceToByteSlice(results)
	a.cache.Add(url, cacheEntry)
	return nil
}
