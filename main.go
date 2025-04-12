package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Initial Config
	initialConfig := config{
		next:     "",
		previous: "",
	}

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
					err := cliCommand.callback(&initialConfig)
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
	callback    func(config *config) error
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

// Map Command
func commandMap(config *config) error {
	fullLocationAreaUrl := BaseUrl + EndpointLocationArea

	var urlToCall string
	if config.previous == "" && config.next == "" {
		// Initial Condition
		urlToCall = fullLocationAreaUrl
	} else if config.next != "" {
		/// After first page has been displayed
		urlToCall = config.next
	}
	// TODO: - implement for previous

	// Get
	res, err := http.Get(urlToCall)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Decode Json
	decoder := json.NewDecoder(res.Body)
	var responseData locationAreaResponse
	if err := decoder.Decode(&responseData); err != nil {
		return err
	}

	// Update next and previous
	config.next = responseData.Next
	config.previous = responseData.Previous

	for _, locationArea := range responseData.Results {
		fmt.Println(locationArea.Name)
	}

	return nil
}

// Exit Command
func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Help Command
func commandHelp(config *config) error {
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
