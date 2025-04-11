package main

import (
	"fmt"
	"strings"
)

func main() {
	finalResult := cleanInput("Hello World")
	fmt.Println(finalResult)
}

func cleanInput(text string) []string {
	fields := strings.Fields(text)

	result := []string{}
	for _, word := range fields {
		result = append(result, strings.ToLower(word))
	}

	return result
}
