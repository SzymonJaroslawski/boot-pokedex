package pokeapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/SzymonJaroslawski/pokedex-go/pokeapi/pokemons"
)

type PokeApiLocationAreaResponse struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
	Count int `json:"count"`
}

type PokeApiLocationAreaWithPokemonsResponse struct {
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name                 string `json:"name"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			Rate int `json:"rate"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
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
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			EncounterDetails []struct {
				Method struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				MaxLevel        int   `json:"max_level"`
				MinLevel        int   `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
}

func GetLocationArea(url string, client *http.Client, cache *Cache) (PokeApiLocationAreaResponse, error) {
	var response PokeApiLocationAreaResponse

	data, ok := cache.Get(url)
	if ok {
		decoder := json.NewDecoder(bytes.NewBuffer(data))
		if err := decoder.Decode(&response); err != nil {
			return PokeApiLocationAreaResponse{}, err
		}
		return response, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokeApiLocationAreaResponse{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return PokeApiLocationAreaResponse{}, err
	}
	defer res.Body.Close()

	byt, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeApiLocationAreaResponse{}, err
	}

	err = json.Unmarshal(byt, &response)
	if err != nil {
		return PokeApiLocationAreaResponse{}, err
	}

	cache.Add(url, byt)

	return response, nil
}

func GetPokemon(url, pokemon string, client *http.Client, cache *Cache) (pokemons.Pokemon, error) {
	var response pokemons.Pokemon

	url = url + "/" + pokemon

	data, ok := cache.Get(url)
	if ok {
		decoder := json.NewDecoder(bytes.NewBuffer(data))
		if err := decoder.Decode(&response); err != nil {
			return pokemons.Pokemon{}, err
		}
		return response, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return pokemons.Pokemon{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return pokemons.Pokemon{}, err
	}
	defer res.Body.Close()

	byt, err := io.ReadAll(res.Body)
	if err != nil {
		return pokemons.Pokemon{}, err
	}

	err = json.Unmarshal(byt, &response)
	if err != nil {
		return pokemons.Pokemon{}, err
	}

	cache.Add(url, byt)

	return response, nil
}

func GetLocationAreaWithPokemons(url, location string, client *http.Client, cache *Cache) (PokeApiLocationAreaWithPokemonsResponse, error) {
	var response PokeApiLocationAreaWithPokemonsResponse

	url = url + "/" + location

	data, ok := cache.Get(url)
	if ok {
		decoder := json.NewDecoder(bytes.NewBuffer(data))
		if err := decoder.Decode(&response); err != nil {
			return PokeApiLocationAreaWithPokemonsResponse{}, err
		}
		return response, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return PokeApiLocationAreaWithPokemonsResponse{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return PokeApiLocationAreaWithPokemonsResponse{}, err
	}
	defer res.Body.Close()

	byt, err := io.ReadAll(res.Body)
	if err != nil {
		return PokeApiLocationAreaWithPokemonsResponse{}, err
	}

	err = json.Unmarshal(byt, &response)
	if err != nil {
		return PokeApiLocationAreaWithPokemonsResponse{}, err
	}

	cache.Add(url, byt)

	return response, nil
}
