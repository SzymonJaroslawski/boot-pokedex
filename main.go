package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SzymonJaroslawski/pokedex-go/commands"
	"github.com/SzymonJaroslawski/pokedex-go/pokeapi"
	"github.com/SzymonJaroslawski/pokedex-go/pokeapi/pokemons"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	interval := time.Second * 5
	cfg := commands.Config{
		Cmds:            commands.GetCommands(),
		PokeApiClient:   &http.Client{},
		Cache:           pokeapi.NewCache(interval),
		CatchedPokemons: make(map[string]pokemons.Pokemon),
	}

	fmt.Print("Welcome to pokedex-go: \n use 'help' to display avaible commands\n")

	for {
		fmt.Print("pokedex-go > ")
		if scanner.Scan() {
			commandInput := scanner.Text()
			commandInputWithParams := strings.Split(commandInput, " ")
			if cmd, ok := cfg.Cmds[commandInputWithParams[0]]; ok {
				fmt.Print("\n")
				err := cmd.Callback(&cfg, commandInputWithParams[1:]...)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}
