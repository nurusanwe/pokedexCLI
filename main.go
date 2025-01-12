package main

import (
	"time"

	"github.com/nurusanwe/pokedexcli/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 10*time.Minute)
	cfg := &config{
		pokeapiClient: pokeClient,
		caughtPokemon: map[string]pokeapi.PokemonDetails{},
	}

	startRepl(cfg)
}
