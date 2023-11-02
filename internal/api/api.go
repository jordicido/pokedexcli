package Api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	Cache "github.com/jordicido/pokedexcli/internal/cache"
)

type LocationArea struct {
	Name string
	Url  string
}

type Pokemon struct {
	Name   string
	Height int
	Weight int
	Stats  []struct {
		BaseStat int
		Name     string
	}
	Types []struct {
		Name string
	}
}

var Pokedex map[string]Pokemon

var cache *Cache.Cache

func GetLocations(link *string) ([]LocationArea, *string, *string) {
	var locations []LocationArea
	type response struct {
		Count    int     `json:"count"`
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}

	if cache == nil {
		cache = Cache.NewCache(5 * time.Second)
	}
	var res *http.Response
	var body []byte
	var err error
	if link == nil {
		res, err = http.Get("https://pokeapi.co/api/v2/location-area/")
		if err != nil {
			log.Fatal(err)
		}
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			log.Fatal(err)
		}
		cache.Add("https://pokeapi.co/api/v2/location-area/", body)
	} else {
		cachedData, ok := cache.Get(*link)
		if ok {
			body = cachedData
		} else {
			res, err = http.Get(*link)
			if err != nil {
				log.Fatal(err)
			}
			body, err = io.ReadAll(res.Body)
			res.Body.Close()
			if res.StatusCode > 299 {
				log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
			}
			if err != nil {
				log.Fatal(err)
			}
			cache.Add(*link, body)
		}

	}

	apiResponse := &response{}
	err = json.Unmarshal(body, apiResponse)
	if err != nil {
		fmt.Println(err)
	}
	for _, location := range apiResponse.Results {
		locations = append(locations, LocationArea{Name: location.Name, Url: location.URL})
	}
	return locations, apiResponse.Next, apiResponse.Previous
}

func GetPokemonsInArea(area string) []string {
	var pokemons []string
	url := "https://pokeapi.co/api/v2/location-area/"
	type response struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}
	var body []byte
	var err error
	cachedData, ok := cache.Get(url + area)
	if ok {
		body = cachedData
	} else {
		res, err := http.Get(url + area)
		if err != nil {
			log.Fatal(err)
		}
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		}
		if err != nil {
			log.Fatal(err)
		}
		cache.Add(url+area, body)
	}

	apiResponse := &response{}
	err = json.Unmarshal(body, apiResponse)
	if err != nil {
		fmt.Println(err)
	}
	for _, pokemon := range apiResponse.PokemonEncounters {
		pokemons = append(pokemons, pokemon.Pokemon.Name)
	}
	return pokemons
}

func CatchPokemon(name string) bool {
	if Pokedex == nil {
		Pokedex = make(map[string]Pokemon)
	}
	url := "https://pokeapi.co/api/v2/pokemon/"
	type response struct {
		BaseExperience int    `json:"base_experience"`
		Height         int    `json:"height"`
		Name           string `json:"name"`
		Stats          []struct {
			BaseStat int `json:"base_stat"`
			Effort   int `json:"effort"`
			Stat     struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"stat"`
		} `json:"stats"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
		Weight int `json:"weight"`
	}

	res, err := http.Get(url + name)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	apiResponse := &response{}
	err = json.Unmarshal(body, apiResponse)
	if err != nil {
		fmt.Println(err)
	}
	catchOption := rand.Intn(apiResponse.BaseExperience)
	result := catchOption < 75
	if result {
		pokemon := Pokemon{
			Name:   apiResponse.Name,
			Height: apiResponse.Height,
			Weight: apiResponse.Weight,
		}
		for _, stat := range apiResponse.Stats {
			pokemon.Stats = append(pokemon.Stats, struct {
				BaseStat int
				Name     string
			}{BaseStat: stat.BaseStat, Name: stat.Stat.Name})
		}
		for _, pokemonType := range apiResponse.Types {
			pokemon.Types = append(pokemon.Types, struct{ Name string }{Name: pokemonType.Type.Name})
		}
		Pokedex[name] = pokemon
	}

	return result
}

func InspectPokemon(name string) (Pokemon, bool) {
	pokemon, ok := Pokedex[name]
	return pokemon, ok
}

func GetPokedex() []string {
	var pokemonList []string
	for _, pokemon := range Pokedex {
		pokemonList = append(pokemonList, pokemon.Name)
	}
	return pokemonList
}
