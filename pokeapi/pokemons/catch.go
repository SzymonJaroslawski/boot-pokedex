package pokemons

import "fmt"

func CommandInspect(pokemons map[string]Pokemon, pokemon_name string) string {
	msg := ""
	if pokemon, ok := pokemons[pokemon_name]; ok {
		msg += fmt.Sprintf("Name: %s\n", pokemon.Name)
		msg += fmt.Sprintf("Height: %d\n", pokemon.Height)
		msg += fmt.Sprintf("Weight: %d\n", pokemon.Weight)
		msg += fmt.Sprintln("Stats:")
		for _, stat := range pokemon.Stats {
			msg += fmt.Sprintf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		msg += fmt.Sprintln("Types:")
		for _, typ := range pokemon.Types {
			msg += fmt.Sprintf(" - %s\n", typ.Type.Name)
		}

		return msg
	}

	msg = "you have not caught that pokemon"
	return msg
}
