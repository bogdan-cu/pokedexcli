package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	PrevUrl string
	NextUrl string
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationList struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Location struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon        Pokemon `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetLocationArea(config *Config, forward bool) ([]string, error) {
	var url string
	if forward {
		if config.NextUrl == "" {
			return nil, fmt.Errorf("you're on the last page")
		}
		url = config.NextUrl
	} else {
		if config.PrevUrl == "" {
			return nil, fmt.Errorf("you're on the first page")
		}
		url = config.PrevUrl
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var locations LocationList
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&locations); err != nil {
		return nil, err
	}

	if forward {
		config.PrevUrl = config.NextUrl
		if locations.Next == nil {
			config.NextUrl = ""
		}
		config.NextUrl = *locations.Next
	} else {
		config.NextUrl = config.PrevUrl
		if locations.Previous == nil {
			config.PrevUrl = ""
		}
		config.PrevUrl = *locations.Next
	}

	var areas []string
	for _, area := range locations.Results {
		areas = append(areas, area.Name)
	}
	return areas, nil
}

func GetLocalPokemon(locationUrl string) ([]string, error) {
	res, err := http.Get(locationUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var payload Location
	err = decoder.Decode(&payload)
	if err != nil {
		return nil, err
	}

	var pokemon []string
	for _, encounter := range payload.PokemonEncounters {
		pokemon = append(pokemon, encounter.Pokemon.Name)
	}
	return pokemon, nil
}
