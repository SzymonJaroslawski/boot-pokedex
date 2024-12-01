package commands

import (
	"fmt"
	"os"
)

type CliCommand struct {
	Callback    func(map[string]CliCommand) error
	Name        string
	Description string
}

func GetCommands() map[string]CliCommand {
	commandMap := map[string]CliCommand{
		"help": {
			Name:        "help",
			Description: "Displays help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit pokedex-go",
			Callback:    commandExit,
		},
	}
	return commandMap
}

func commandHelp(cmds map[string]CliCommand) error {
	helpMsg := ""
	for _, cmd := range cmds {
		helpMsg += fmt.Sprintf("%s\n%s", cmd.Name, cmd.Description)
	}
	fmt.Print(helpMsg)
	return nil
}

func commandExit(cmds map[string]CliCommand) error {
	os.Exit(0)
	return nil
}
