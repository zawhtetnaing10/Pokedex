package commands

import (
	"github.com/zawhtetnaing10/Pokedex/internal/caughtpokemon"
	"github.com/zawhtetnaing10/Pokedex/internal/networkresponses"
	"github.com/zawhtetnaing10/Pokedex/internal/pokecache"
)

type CliCommand struct {
	Name        string
	Description string
	Callback    func(config *networkresponses.Config, cache *pokecache.Cache, firstArg string, pokedex *caughtpokemon.Pokedex) error
}
