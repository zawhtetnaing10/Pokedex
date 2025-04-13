package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zawhtetnaing10/Pokedex/internal/pokecache"
)

// Supported commands
var supportedCommands map[string]cliCommand

func init() {
	supportedCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays location areas in Pokemon World",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous location areas in Pokemon World",
			callback:    commandMapBack,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Initial Config
	initialConfig := config{
		next:     "",
		previous: "",
	}

	// Cache
	pokecache := pokecache.NewCache(5 * time.Second)

	for {
		fmt.Print("Pokedex >")

		var input string
		if scanner.Scan() {
			input = scanner.Text()

			cleanedWords := cleanInput(input)

			if len(cleanedWords) > 0 {
				inputCommand := cleanedWords[0]

				cliCommand, ok := supportedCommands[inputCommand]

				if ok {
					err := cliCommand.callback(&initialConfig, pokecache)
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

type config struct {
	previous string
	next     string
}

type cliCommand struct {
	name        string
	description string
	callback    func(config *config, cache *pokecache.Cache) error
}

type locationAreaResponse struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []locationAreaData `json:"results"`
}

type locationAreaData struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// Map back command
func commandMapBack(config *config, cache *pokecache.Cache) error {
	if config.previous == "" {
		fmt.Println("You are on the first page")
		return nil
	} else {
		// Make Api Call
		responseData, err := getResponseFromRepo(config.previous, cache)
		if err != nil {
			return err
		}

		// Update next and previous
		config.next = responseData.Next
		config.previous = responseData.Previous

		// Print out the results
		for _, locationArea := range responseData.Results {
			fmt.Println(locationArea.Name)
		}

		return nil
	}
}

// Map Command
func commandMap(config *config, cache *pokecache.Cache) error {
	fullLocationAreaUrl := BaseUrl + EndpointLocationArea

	var urlToCall string
	if config.previous == "" && config.next == "" {
		// Initial Condition
		urlToCall = fullLocationAreaUrl
	} else if config.next != "" {
		/// After first page has been displayed
		urlToCall = config.next
	}

	// Get from cache or from api.
	responseData, err := getResponseFromRepo(urlToCall, cache)
	if err != nil {
		return err
	}

	// Update next and previous
	config.next = responseData.Next
	config.previous = responseData.Previous

	// Print out the results
	for _, locationArea := range responseData.Results {
		fmt.Println(locationArea.Name)
	}

	return nil
}

// Checks if the resource is in the cache, if not make an api call and update the cache
func getResponseFromRepo(url string, cache *pokecache.Cache) (locationAreaResponse, error) {
	// Accesses Cache
	responseDataFromCache, found, err := findResponseFromCache(url, cache)
	if err != nil {
		return locationAreaResponse{}, fmt.Errorf("cached data may be corrupted. it cannot be parsed %w", err)
	}

	if found {
		// Cache Exists
		return responseDataFromCache, nil
	} else {
		// Make Api Call
		responseDataFromApi, err := makeGetApiCall[locationAreaResponse](url)
		if err != nil {
			return locationAreaResponse{}, fmt.Errorf("error making api call %w", err)
		}

		// Update cache
		cacheBytes, err := json.Marshal(responseDataFromApi)
		if err != nil {
			return locationAreaResponse{}, fmt.Errorf("error marshalling network response %w", err)
		}
		cache.Add(url, cacheBytes)

		// Returns the api response
		return responseDataFromApi, nil
	}
}

// Checks if the resource already exists with the url
func findResponseFromCache(url string, cache *pokecache.Cache) (locationAreaResponse, bool, error) {
	cachedValue, found := cache.Get(url)
	if found {
		var cachedResponse locationAreaResponse
		err := json.Unmarshal(cachedValue, &cachedResponse)
		if err != nil {
			return locationAreaResponse{}, false, fmt.Errorf("error unmarshalling cached data %w", err)
		} else {
			return cachedResponse, found, nil
		}
	}
	return locationAreaResponse{}, false, nil
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
func commandExit(config *config, cache *pokecache.Cache) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Help Command
func commandHelp(config *config, cache *pokecache.Cache) error {
	helpMessage := "Welcome to the Pokedex!\nUsage:\n\n"

	for key, value := range supportedCommands {
		messageToAdd := fmt.Sprintf("%s: %s\n", key, value.description)
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
