package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func commandExit(conf *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(conf *config, args ...string) error {
	locationsResp, err := conf.pokeapiClient.ListLocations(conf.nextLocationsURL)
	if err != nil {
		return err
	}

	conf.nextLocationsURL = locationsResp.Next
	conf.prevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapBack(conf *config, args ...string) error {
	if conf.prevLocationsURL == nil {
		return errors.New("you're on the first page")
	}

	locationResp, err := conf.pokeapiClient.ListLocations(conf.prevLocationsURL)
	if err != nil {
		return err
	}

	conf.nextLocationsURL = locationResp.Next
	conf.prevLocationsURL = locationResp.Previous

	for _, loc := range locationResp.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("you must provide a unique location name")
	}

	name := args[0]
	location, err := cfg.pokeapiClient.ListExplore(name)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", location.Name)
	fmt.Println("Found Pokemon: ")
	for _, enc := range location.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}
	return nil
}

func commandCatch(conf *config, args ...string) error {
	if len(args) > 1 {
		return errors.New("you can only catch one Pokemon at a time")
	}

	if len(args) == 0 {
		if len(conf.caughtPokemon) == 0 {
			fmt.Println("No Pokemon caught yet.")
			return nil
		}
		fmt.Println("Caught Pokemon:")
		for name := range conf.caughtPokemon {
			fmt.Println(" -", name)
		}
		return nil
	}

	pokemonName := strings.ToLower(args[0])

	if _, exists := conf.caughtPokemon[pokemonName]; exists {
		fmt.Printf("%s has already been caught!\n", pokemonName)
		return nil
	}

	pokemonDetails, err := conf.pokeapiClient.FetchPokemonDetails(pokemonName)
	if err != nil {
		return err
	}

	// Define a fixed threshold
	catchThreshold := pokemonDetails.BaseExperience / 2
	// Generate one random number, no Seed setup needed anymore since go 1.20
	catchAttempt := rand.Intn(pokemonDetails.BaseExperience)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonDetails.Name)
	if catchAttempt <= catchThreshold {
		fmt.Printf("%s was caught!\n", pokemonDetails.Name)
		conf.caughtPokemon[pokemonDetails.Name] = pokemonDetails
	} else {
		fmt.Printf("%s escaped!\n", pokemonDetails.Name)
	}

	return nil
}

func commandInspect(conf *config, args ...string) error {

	if len(args) != 1 {
		return errors.New("you must provide a pokemon name")
	}

	name := strings.ToLower(args[0])
	if _, exists := conf.caughtPokemon[name]; !exists {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	pokemon := conf.caughtPokemon[name]
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(conf *config, args ...string) error {
	fmt.Println("Your Pokedex:")
	for name := range conf.caughtPokemon {
		fmt.Println(" -", name)
	}
	return nil
}
