package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	Api "github.com/jordicido/pokedexcli/internal/api"
)

type cliCommand struct {
	name        string
	description string
	callback    func(interface{}) error
}

var config struct {
	Next     *string
	Previous *string
}

var commands map[string]cliCommand

func commandHelp(interface{}) error {
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	for name := range commands {
		fmt.Printf("%s: %s \n", commands[name].name, commands[name].description)
	}
	return nil
}

func commandExit(interface{}) error {
	os.Exit(0)
	return nil
}

func commandMap(interface{}) error {
	var locations []Api.LocationArea
	locations, config.Next, config.Previous = Api.GetLocations(config.Next)
	for _, location := range locations {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapB(interface{}) error {
	var locations []Api.LocationArea
	locations, config.Next, config.Previous = Api.GetLocations(config.Previous)
	for _, location := range locations {
		fmt.Println(location.Name)
	}

	return nil
}

func commandExplore(area interface{}) error {
	fmt.Printf("Exploring %s...\n", area)
	fmt.Println("Found Pokemon:")

	pokemons := Api.GetPokemonsInArea(area.(string))
	for _, pokemon := range pokemons {
		fmt.Printf("- %s\n", pokemon)
	}

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
		input := strings.Split(scanner.Text(), " ")
		command, ok := commands[input[0]]
		if ok {
			var param interface{}
			if len(input) > 1 {
				param = input[1]
			}
			command.callback(param)
		} else {
			fmt.Println("This command is not valid")
		}
	}
}
