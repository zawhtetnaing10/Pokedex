package caughtpokemon

import (
	"sync"

	"github.com/zawhtetnaing10/Pokedex/internal/networkresponses"
)

type Pokedex struct {
	CaughtPokemon map[string]networkresponses.Pokemon
	Mu            sync.RWMutex
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		CaughtPokemon: make(map[string]networkresponses.Pokemon),
	}
}

func (p *Pokedex) Add(name string, pokemon networkresponses.Pokemon) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.CaughtPokemon[name] = pokemon
}

func (p *Pokedex) Get(name string) networkresponses.Pokemon {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	return p.CaughtPokemon[name]
}
