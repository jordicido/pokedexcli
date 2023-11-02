package Api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	Cache "github.com/jordicido/pokedexcli/internal/cache"
)

type LocationArea struct {
	Name string
	Url  string
}

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
			Pokemon struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"pokemon"`
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
