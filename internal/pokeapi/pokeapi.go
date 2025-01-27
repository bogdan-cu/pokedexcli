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

	areas := []string{}
	for _, area := range locations.Results {
		areas = append(areas, area.Name)
	}
	return areas, nil
}

/*
state 1 - prev empty, next value (0-20)
forward - using next
state 2 - prev value(0-20), next value(20-40)
back - using prev
state 3 - prev empty, next value(0-20)
*/
