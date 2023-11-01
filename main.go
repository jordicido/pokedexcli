package main

import (
	"bufio"
	"fmt"
	"os"

	Api "github.com/jordicido/pokedexcli/internal/api"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var config struct {
	Next     *string
	Previous *string
}

var commands map[string]cliCommand

func commandHelp() error {
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	for name := range commands {
		fmt.Printf("%s: %s \n", commands[name].name, commands[name].description)
	}
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func commandMap() error {
	var locations []string
	locations, config.Next, config.Previous = Api.CallPokeLocationArea(config.Next)
	for _, location := range locations {
		fmt.Println(location)
	}

	return nil
}

func commandMapB() error {
	var locations []string
	locations, config.Next, config.Previous = Api.CallPokeLocationArea(config.Previous)
	for _, location := range locations {
		fmt.Println(location)
	}

	return nil
}

func commandExplore() error {
	//TODO

	return nil
}

func initCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas in the Pokemon world",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "List of all the Pokémon in a given area",
			callback:    commandExplore,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands = initCommands()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		command, ok := commands[input]
		if ok {
			command.callback()
		} else {
			fmt.Println("This command is not valid")
		}
	}
}
