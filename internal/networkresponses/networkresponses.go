package networkresponses

type Config struct {
	Previous string
	Next     string
}

type LocationAreaResponse struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []LocationAreaData `json:"results"`
}

type LocationAreaData struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationAreaWithPokemon struct {
	PokemonEncounters []PokemonEncounters `json:"pokemon_encounters"`
}

type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name           string          `json:"name"`
	Url            string          `json:"url"`
	BaseExperience int             `json:"base_experience"`
	Height         int             `json:"height"`
	Weight         int             `json:"weight"`
	Stats          []StatContainer `json:"stats"`
	Types          []TypeContainer `json:"types"`
}

// Get stat by name
func (p *Pokemon) GetStatByName(statName string) int {
	for _, stat := range p.Stats {
		if stat.Stat.Name == statName {
			return stat.BaseStat
		}
	}
	return 0
}

// Get Types
func (p *Pokemon) GetTypes() []string {
	result := []string{}

	for _, pokemonType := range p.Types {
		result = append(result, pokemonType.Type.Name)
	}

	return result
}

// StatContainer
type StatContainer struct {
	BaseStat int  `json:"base_stat"`
	Stat     Stat `json:"stat"`
}

// Stat
type Stat struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// Type Container
type TypeContainer struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

// Type
type Type struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
