package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zawhtetnaing10/Pokedex/internal/caughtpokemon"
	"github.com/zawhtetnaing10/Pokedex/internal/commands"
	"github.com/zawhtetnaing10/Pokedex/internal/networkresponses"
	"github.com/zawhtetnaing10/Pokedex/internal/pokecache"
)

// Supported commands
var supportedCommands map[string]commands.CliCommand

func init() {
	supportedCommands = map[string]commands.CliCommand{
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"map": {
			Name:        "map",
			Description: "Displays location areas in Pokemon World",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays previous location areas in Pokemon World",
			Callback:    commandMapBack,
		},
		"explore": {
			Name:        "explore",
			Description: "Displays pokemon in the location area",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Catches pokemon",
			Callback:    commandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspects pokemon",
			Callback:    commandInspect,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Initial Config
	initialConfig := networkresponses.Config{
		Next:     "",
		Previous: "",
	}

	// Cache
	pokecache := pokecache.NewCache(5 * time.Second)

	// Pokedex
	pokedex := caughtpokemon.NewPokedex()

	for {
		fmt.Print("Pokedex >")

		var input string
		if scanner.Scan() {
			input = scanner.Text()

			cleanedWords := cleanInput(input)

			if len(cleanedWords) > 0 {
				// Command
				inputCommand := cleanedWords[0]

				// First arg exists
				var firstArg string
				if len(cleanedWords) >= 2 {
					firstArg = cleanedWords[1]
				}

				cliCommand, ok := supportedCommands[inputCommand]

				if ok {
					err := cliCommand.Callback(&initialConfig, pokecache, firstArg, pokedex)
					if err != nil {
						fmt.Println("Error:", err)
					}
				} else {
					fmt.Println("Unknown command")
				}
			}
		}
	}
}

// Inspect
func commandInspect(config *networkresponses.Config, cache *pokecache.Cache, firstArg string, pokedex *caughtpokemon.Pokedex) error {
	caughtPokemon, ok := pokedex.CaughtPokemon[firstArg]
	if !ok {
		// Pokemon hasn't been caught yet
		fmt.Printf("Cannot find %v in pokedex...\n", firstArg)
		return nil
	} else {
		// Pokemon is already caught. Print out the stats
		// Name, Height, Weight
		fmt.Printf("Name: %v\n", caughtPokemon.Name)
		fmt.Printf("Height: %v\n", caughtPokemon.Height)
		fmt.Printf("Weight: %v\n", caughtPokemon.Weight)

		// Stats
		fmt.Println("Stats:")
		fmt.Printf("  -hp: %v\n", caughtPokemon.GetStatByName("hp"))
		fmt.Printf("  -attack: %v\n", caughtPokemon.GetStatByName("attack"))
		fmt.Printf("  -defense: %v\n", caughtPokemon.GetStatByName("defense"))
		fmt.Printf("  -special-attack: %v\n", caughtPokemon.GetStatByName("special-attack"))
		fmt.Printf("  -special-defense: %v\n", caughtPokemon.GetStatByName("special-defense"))
		fmt.Printf("  -speed: %v\n", caughtPokemon.GetStatByName("speed"))

		// Types
		fmt.Println("Types:")
		for _, pokemonTypeName := range caughtPokemon.GetTypes() {
			fmt.Printf("  - %v\n", pokemonTypeName)
		}

		return nil
	}
}

// Catch command
func commandCatch(config *networkresponses.Config, cache *pokecache.Cache, firstArg string, pokedex *caughtpokemon.Pokedex) error {
	fmt.Printf("Throwing a Pokeball at %v...\n", firstArg)

	pokemonUrl := BaseUrl + EndpointPokemon + "/" + firstArg

	pokemon, err := getResponseFromRepo[networkresponses.Pokemon](pokemonUrl, cache)
	if err != nil {
		return err
	}

	if calculateChanceToCatch(pokemon.BaseExperience) {
		// Pokemon caught
		pokedex.Add(pokemon.Name, pokemon)
		fmt.Printf("%v was caught!\n", pokemon.Name)
	} else {
		// Pokemon Escaped
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}
	return nil
}

// Calculates the chance for success
func calculateChanceToCatch(baseExperience int) bool {
	randomValue := rand.Intn(MaxCatchThreshold)

	return randomValue >= baseExperience && randomValue < MaxCatchThreshold
}

// Explore command
func commandExplore(config *networkresponses.Config, cache *pokecache.Cache, firstArg string, pokedex *caughtpokemon.Pokedex) error {
	fmt.Printf("Exploring %v...\n", firstArg)
	fmt.Println("Found Pokemon:")

	fullPokemonListUrl := BaseUrl + EndpointLocationArea + "/" + firstArg

	response, err := getResponseFromRepo[networkresponses.LocationAreaWithPokemon](fullPokemonListUrl, cache)
	if err != nil {
		return err
	}

	for _, pokemonDetails := range response.PokemonEncounters {
		fmt.Printf("%v\n", pokemonDetails.Pokemon.Name)
	}

	return nil
}

// Map back command
func commandMapBack(config *networkresponses.Config, cache *pokecache.Cache, firstArg string, pokedex *caughtpokemon.Pokedex) error {
	if config.Previous == "" {
		fmt.Println("You are on the first page")
		return nil
	} else {
		// Make Api Call
		responseData, err := getResponseFromRepo[networkresponses.LocationAreaResponse](config.Previous, cache)
		if err != nil {
			return err
		}

		// Update next and previous
		config.Next = responseData.Next
		config.Previous = responseData.Previous

		// Print out the results
		for _, locationArea := range responseData.Results {
			fmt.Println(locationArea.Name)
		}

		return nil
	}
}

// Map Command
func commandMap(config *networkresponses.Config, cache *pokecache.Cache, firstArg string, pokedex *caughtpokemon.Pokedex) error {
	fullLocationAreaUrl := BaseUrl + EndpointLocationArea

	var urlToCall string
	if config.Previous == "" && config.Next == "" {
		// Initial Condition
		urlToCall = fullLocationAreaUrl
	} else if config.Next != "" {
		/// After first page has been displayed
		urlToCall = config.Next
	}

	// Get from cache or from api.
	responseData, err := getResponseFromRepo[networkresponses.LocationAreaResponse](urlToCall, cache)
	if err != nil {
		return err
	}

	// Update next and previous
	config.Next = responseData.Next
	config.Previous = responseData.Previous

	// Print out the results
	for _, locationArea := range responseData.Results {
		fmt.Println(locationArea.Name)
	}

	return nil
}

// Checks if the resource is in the cache, if not make an api call and update the cache
func getResponseFromRepo[T any](url string, cache *pokecache.Cache) (T, error) {
	// Accesses Cache
	responseDataFromCache, found, err := findResponseFromCache[T](url, cache)
	if err != nil {
		return *new(T), fmt.Errorf("cached data may be corrupted. it cannot be parsed %w", err)
	}

	if found {
		// Cache Exists
		return responseDataFromCache, nil
	} else {
		// Make Api Call
		responseDataFromApi, err := makeGetApiCall[T](url)
		if err != nil {
			return *new(T), fmt.Errorf("error making api call %w", err)
		}

		// Update cache
		cacheBytes, err := json.Marshal(responseDataFromApi)
		if err != nil {
			return *new(T), fmt.Errorf("error marshalling network response %w", err)
		}
		cache.Add(url, cacheBytes)

		// Returns the api response
		return responseDataFromApi, nil
	}
}

// Checks if the resource already exists with the url
func findResponseFromCache[T any](url string, cache *pokecache.Cache) (T, bool, error) {
	cachedValue, found := cache.Get(url)
	if found {
		var cachedResponse T
		err := json.Unmarshal(cachedValue, &cachedResponse)
		if err != nil {
			return *new(T), false, fmt.Errorf("error unmarshalling cached data %w", err)
		} else {
			return cachedResponse, found, nil
		}
	}
	return *new(T), false, nil
}

// Generic function to make Get API Call
func makeGetApiCall[T any](urlToCall string) (T, error) {
	// Make Get Api Call
	res, err := http.Get(urlToCall)
	if err != nil {
		return *new(T), fmt.Errorf("failed to fetch api %w", err)
	}
	defer res.Body.Close()

	// Decode Json
	decoder := json.NewDecoder(res.Body)
	var responseData T
	if err := decoder.Decode(&responseData); err != nil {
		return *new(T), fmt.Errorf("failed to decode json %w", err)
	}

	return responseData, nil
}

// Exit Command
func commandExit(config *networkresponses.Config, cache *pokecache.Cache, firstArg string, pokedex *caughtpokemon.Pokedex) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Help Command
func commandHelp(config *networkresponses.Config, cache *pokecache.Cache, firstArg string, pokedex *caughtpokemon.Pokedex) error {
	helpMessage := "Welcome to the Pokedex!\nUsage:\n\n"

	for key, value := range supportedCommands {
		messageToAdd := fmt.Sprintf("%s: %s\n", key, value.Description)
		helpMessage += messageToAdd
	}
	fmt.Println(helpMessage)
	return nil
}

func cleanInput(text string) []string {
	fields := strings.Fields(text)

	result := []string{}
	for _, word := range fields {
		result = append(result, strings.ToLower(word))
	}

	return result
}
