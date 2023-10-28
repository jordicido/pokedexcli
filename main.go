package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	Api "github.com/jordicido/pokedexcli/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func refreshConfig(advance bool) {
	if advance {
		if config.Next > 0 {
			config.Previous += 20
		}
		config.Next += 20
	} else {
		if config.Next == 0 {
			fmt.Println("This command is not valid")
			return
		} else if config.Previous > 0 {
			config.Previous -= 20
		}
		config.Next -= 20
	}
}

var config struct {
	Next     int
	Previous int
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
	nextOffset := strconv.Itoa(config.Next)
	locations := Api.CallPokeLocationArea(nextOffset)
	for _, location := range locations {
		fmt.Println(location)
	}
	refreshConfig(true)

	return nil
}

func commandMapB() error {
	previousOffset := strconv.Itoa(config.Previous)
	locations := Api.CallPokeLocationArea(previousOffset)
	for _, location := range locations {
		fmt.Println(location)
	}
	refreshConfig(false)

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
