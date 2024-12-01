package commands

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"

	"github.com/SzymonJaroslawski/pokedex-go/pokeapi"
	"github.com/SzymonJaroslawski/pokedex-go/pokeapi/pokemons"
)

const (
	basePokeApiLocationAreURL          = "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20"
	basePokeApiLocationWithPokemonsURL = "https://pokeapi.co/api/v2/location-area"
	basePokeApiCatchURL                = "https://pokeapi.co/api/v2/pokemon"
)

type CliCommand struct {
	Callback    func(*Config, ...string) error
	Name        string
	Description string
}

type Config struct {
	Cmds             map[string]CliCommand
	NextLocationArea *string
	PrevLocationArea *string
	PokeApiClient    *http.Client
	Cache            *pokeapi.Cache
	CatchedPokemons  map[string]pokemons.Pokemon
}

func GetCommands() map[string]CliCommand {
	cmds := map[string]CliCommand{
		"help": {
			Name:        "help",
			Description: "Displays help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exits program",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "Display 20 locations, each use displays next 20, use pmap to go back 20 locations",
			Callback:    commandMap,
		},
		"pmap": {
			Name:        "pmap",
			Description: "Display previous 20 locations, use map to go 20 locations forward \n Returns error if map wasn't used at least one time",
			Callback:    commandPmap,
		},
		"explore": {
			Name:        "explore <location>",
			Description: "Displays all pokemon that can be found in location",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch <pokemon_name>",
			Description: "Use to try catching pokemon, catch change depends on pokemons base experience",
			Callback:    commandCatch,
		},
		"inspect": {
			Name:        "inspect <pokemon_name>",
			Description: "Use to inspect pokemons stats, you must catch it first",
			Callback:    commandInspect,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Use to display all caught pokemons",
			Callback:    commandPokedex,
		},
	}
	return cmds
}

func commandHelp(cfg *Config, params ...string) error {
	msg := ""
	for _, cmd := range cfg.Cmds {
		msg += fmt.Sprintf("%s\n%s\n\n", cmd.Name, cmd.Description)
	}
	fmt.Print(msg)
	return nil
}

func commandExit(cfg *Config, params ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(cfg *Config, params ...string) error {
	URL := basePokeApiLocationAreURL
	if cfg.NextLocationArea != nil {
		URL = *cfg.NextLocationArea
	}

	res, err := pokeapi.GetLocationArea(URL, cfg.PokeApiClient, cfg.Cache)
	if err != nil {
		return err
	}

	cfg.NextLocationArea = &res.Next
	cfg.PrevLocationArea = &res.Previous

	for _, location := range res.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func commandPmap(cfg *Config, params ...string) error {
	if cfg.PrevLocationArea == nil {
		return fmt.Errorf("No locations to show, use command map at least once")
	}

	if *cfg.PrevLocationArea == "" {
		return fmt.Errorf("Can't go back already at the start")
	}

	res, err := pokeapi.GetLocationArea(*cfg.PrevLocationArea, cfg.PokeApiClient, cfg.Cache)
	if err != nil {
		return err
	}

	cfg.NextLocationArea = &res.Next
	cfg.PrevLocationArea = &res.Previous

	for _, location := range res.Results {
		fmt.Printf("%s\n", location.Name)
	}

	return nil
}

func commandExplore(cfg *Config, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("please specify location to explore (use 'help' to see manual)")
	}

	res, err := pokeapi.GetLocationAreaWithPokemons(basePokeApiLocationWithPokemonsURL, params[0], cfg.PokeApiClient, cfg.Cache)
	if err != nil {
		return err
	}

	fmt.Printf("exploring %s...\n", params[0])

	fmt.Printf("found pokemons: \n")

	for _, pokemon := range res.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *Config, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("please provide pokemon name to catch\n")
	}

	res, err := pokeapi.GetPokemon(basePokeApiCatchURL, params[0], cfg.PokeApiClient, cfg.Cache)
	if err != nil {
		return err
	}

	fmt.Printf("throwing a Pokeball at %s\n", params[0])

	if rand.IntN(res.BaseExperience/50) == (res.BaseExperience/50 - 1) {
		fmt.Printf("%s was caught!\n", res.Name)
		cfg.CatchedPokemons[res.Name] = res
		return nil
	}

	fmt.Printf("%s escaped!\n", res.Name)
	return nil
}

func commandInspect(cfg *Config, params ...string) error {
	if len(params) == 0 {
		return fmt.Errorf("please provide pokemon name to inspect\n")
	}

	msg := pokemons.CommandInspect(cfg.CatchedPokemons, params[0])

	fmt.Println(msg)
	return nil
}

func commandPokedex(cfg *Config, params ...string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.CatchedPokemons {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}
