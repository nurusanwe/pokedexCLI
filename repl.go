package main

import (
	"fmt"
	"strings"

	"github.com/nurusanwe/pokedexcli/internal/pokeapi"
	"github.com/peterh/liner"
)

type config struct {
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
	caughtPokemon    map[string]pokeapi.PokemonDetails
}

func startRepl(cfg *config) {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetCompleter(func(line string) (c []string) {
		for _, n := range getCommands() {
			if strings.HasPrefix(n.name, strings.ToLower(line)) {
				c = append(c, n.name)
			}
		}
		return
	})

	fmt.Println("Welcome to the Pokedex CLI!")
	fmt.Println("Type 'help' to see available commands.")

	for {
		cmd, err := line.Prompt("Pokedex > ")
		if err != nil {
			if err == liner.ErrPromptAborted {
				fmt.Println("Aborted")
				return
			}
			fmt.Println("Error reading line:", err)
			continue
		}

		line.AppendHistory(cmd)
		words := cleanInput(cmd)

		if len(words) == 0 {
			continue
		}

		firstWord := words[0]
		command, found := getCommands()[firstWord]
		if !found {
			fmt.Printf("Unknown command: %s\n", firstWord)
			continue
		}

		err = command.callback(cfg, words[1:]...)
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	for i, word := range words {
		words[i] = strings.ToLower(strings.TrimSpace(word))
	}
	return words
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getCommands() map[string]cliCommand {
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
			description: "Display the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore <location>",
			description: "Display the pokemon in a location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon>",
			description: "Attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon>",
			description: "Provide details on a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Provide the list of caught pokemon",
			callback:    commandPokedex,
		},
	}
}
