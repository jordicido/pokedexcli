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

var cache *Cache.Cache

func CallPokeLocationArea(link *string) ([]string, *string, *string) {

	var locations []string
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
		locations = append(locations, location.Name)
	}
	return locations, apiResponse.Next, apiResponse.Previous
}
