package main

import (
	"fmt"
	"github.com/bogdan-cu/pokedexcli/internal/pokecache"
	"os"

	"github.com/bogdan-cu/pokedexcli/internal/pokeapi"
)

const locationAreaUrl = "https://pokeapi.co/api/v2/location-area/"
const pokemonUrl = "https://pokeapi.co/api/v2/pokemon/"

type App struct {
	pokedex *pokeapi.Pokedex
	config  *pokeapi.Config
	cache   *pokecache.Cache
}

type CliCommand struct {
	name        string
	description string
	callback    func(*App, string) error
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
		callback:    (*App).commandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Attempt to catch a specific pokemon",
		callback:    (*App).commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "Returns a basic stat sheet for a Pokemon you have in your Pokedex",
		callback:    (*App).commandInspect,
	},
	"pokedex": {
		name:        "pokedex",
		description: "Returns a list of Pokemon you've caught",
		callback:    (*App).commandPokedex,
	},
}

func (a *App) commandExit(_ string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (a *App) commandHelp(_ string) error {
	fmt.Println(`
Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex
map: It displays the names of 20 location areas in the Pokemon world
mapb: It displays the names of the previous 20 location areas in the Pokemon world
explore <map_area>: It returns a list of pokemon found in a given location
catch <pokemon_name>: Attempts to catch a pokemon and add it to the Pokedex
inspect <pokemon_name>: Returns a stat sheet for the Pokemon
pokedex: Returns a list of Pokemon you've caught`)
	return nil
}

func (a *App) commandMap(_ string) error {
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

func (a *App) commandMapb(_ string) error {
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

func (a *App) commandExplore(location string) error {
	locationUrl := locationAreaUrl + location
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

func (a *App) commandCatch(name string) error {
	if _, ok := a.pokedex.Has(name); ok {
		return fmt.Errorf("you already have this one in your pokedex")
	}
	url := pokemonUrl + name
	pokemon, err := pokeapi.GetPokemon(url)
	if err != nil {
		return err
	}
	_ = writeStrings(os.Stdout, "Throwing a Pokeball at "+name+"...")
	if a.pokedex.Catch(pokemon) {
		_ = writeStrings(os.Stdout, name+" was caught!")
	} else {
		_ = writeStrings(os.Stdout, name+" escaped!")
	}

	return nil
}

func (a *App) commandInspect(name string) error {
	pokemon, ok := a.pokedex.Has(name)
	if !ok {
		return fmt.Errorf("you have not caught this Pokemon")
	}
	stats := pokeapi.GetStats(pokemon)
	fmt.Printf("Name: %s\n", stats.Name)
	fmt.Printf("Height: %d\n", stats.Height)
	fmt.Printf("Weight: %d\n", stats.Weight)
	fmt.Println("Stats:")
	for _, stat := range stats.Stats {
		fmt.Printf("\t-%s: %d\n", stat.Name, stat.Val)
	}
	fmt.Println("Types:")
	if err := writeStrings(os.Stdout, stats.Types...); err != nil {
		return err
	}
	return nil
}

func (a *App) commandPokedex(_ string) error {
	fmt.Println("Your Pokedex:")
	if err := writeStrings(os.Stdout, a.pokedex.GetAll()...); err != nil {
		return err
	}
	return nil
}
