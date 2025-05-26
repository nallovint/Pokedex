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
)

type cliCommand struct {
	name        string
	description string
	callback    func(args []string) error
}

func commandExit(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(args []string) error {
	const pageSize = 20
	// Use a static variable to keep track of the current offset between calls
	if mapState.offset == -1 {
		mapState.offset = 0
	}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=%d", mapState.offset, pageSize)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch location areas: %w", err)
	}
	defer resp.Body.Close()
	var result struct {
		Results []struct {
			Name string `json:"name"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	if len(result.Results) == 0 {
		fmt.Println("No more locations to display.")
		return nil
	}
	fmt.Println("Location Areas:")
	for _, loc := range result.Results {
		fmt.Println("-", loc.Name)
	}
	mapState.offset += pageSize
	return nil
}

// mapState keeps track of the current offset for the map command
var mapState = struct{ offset int }{offset: -1}

var pokedex = make(map[string]Pokemon)

type Pokemon struct {
	Name           string
	BaseExperience int `json:"base_experience"`
	Height         int `json:"height"`
	Weight         int `json:"weight"`
	Stats          []struct {
		Name  string
		Value int
	}
	Types []string
}

func commandCatch(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: catch <pokemon>")
		return nil
	}
	pokemonName := args[1]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch pokemon: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("Pokemon '%s' not found.\n", pokemonName)
		return nil
	}
	var pokeData struct {
		Name           string `json:"name"`
		BaseExperience int    `json:"base_experience"`
		Height         int    `json:"height"`
		Weight         int    `json:"weight"`
		Stats          []struct {
			BaseStat int `json:"base_stat"`
			Stat     struct {
				Name string `json:"name"`
			} `json:"stat"`
		} `json:"stats"`
		Types []struct {
			Type struct {
				Name string `json:"name"`
			} `json:"type"`
		} `json:"types"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pokeData); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	catchChance := 100 - pokeData.BaseExperience
	if catchChance < 10 {
		catchChance = 10
	}
	rand.Seed(time.Now().UnixNano())
	roll := rand.Intn(100) + 1 // 1-100
	if roll <= catchChance {
		fmt.Printf("%s was caught!\n", pokemonName)
		stats := make([]struct {
			Name  string
			Value int
		}, len(pokeData.Stats))
		for i, s := range pokeData.Stats {
			stats[i].Name = s.Stat.Name
			stats[i].Value = s.BaseStat
		}
		types := make([]string, len(pokeData.Types))
		for i, t := range pokeData.Types {
			types[i] = t.Type.Name
		}
		pokedex[pokeData.Name] = Pokemon{
			Name:           pokeData.Name,
			BaseExperience: pokeData.BaseExperience,
			Height:         pokeData.Height,
			Weight:         pokeData.Weight,
			Stats:          stats,
			Types:          types,
		}
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	return nil
}

func commandInspect(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: inspect <pokemon>")
		return nil
	}
	pokemonName := args[1]
	poke, ok := pokedex[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", poke.Name)
	fmt.Printf("Height: %d\n", poke.Height)
	fmt.Printf("Weight: %d\n", poke.Weight)
	fmt.Println("Stats:")
	for _, stat := range poke.Stats {
		fmt.Printf("  -%s: %d\n", stat.Name, stat.Value)
	}
	fmt.Println("Types:")
	for _, t := range poke.Types {
		fmt.Printf("  - %s\n", t)
	}
	return nil
}

func commandPokedex(args []string) error {
	if len(pokedex) == 0 {
		fmt.Println("Your Pokedex is empty.")
		return nil
	}
	fmt.Println("Your Pokedex:")
	for name := range pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func main() {
	var commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "exit the pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "open the help",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "display 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"explore": {
			name:        "explore",
			description: "list all Pokémon in a location area (usage: explore <location-area>)",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "try to catch a Pokemon (usage: catch <pokemon>)",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect a caught Pokemon (usage: inspect <pokemon>)",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "list all caught Pokemon",
			callback:    commandPokedex,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		if len(cleanedInput) > 0 {
			command := cleanedInput[0]
			fmt.Printf("You entered %s\n", command)
			if command, ok := commands[command]; ok {
				err := command.callback(cleanedInput)
				if err != nil {
					return
				}
			} else {
				fmt.Println("Command not found")
			}
		}
	}
}

func commandHelp(args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex\nmap: Explore the Pokemon world by listing location areas\nexplore <location-area>: List all Pokémon in a location area\ncatch <pokemon>: Try to catch a Pokémon\ninspect <pokemon>: Inspect a caught Pokémon\npokedex: List all caught Pokémon")
	return nil
}

func cleanInput(text string) []string {
	// hello world -> ["hello", "world"]
	// Charmander Bulbasaur PIKACHU -> ["charmander", "bulbasaur", "pikachu"]
	fields := strings.Fields(strings.TrimSpace(text))
	for i, word := range fields {
		fields[i] = strings.ToLower(word)
	}
	return fields
}

func commandExplore(args []string) error {
	if len(args) < 2 {
		fmt.Println("Usage: explore <location-area>")
		return nil
	}
	locationArea := args[1]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", locationArea)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch location area: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("Location area '%s' not found.\n", locationArea)
		return nil
	}
	var result struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	if len(result.PokemonEncounters) == 0 {
		fmt.Println("No Pokémon found in this location area.")
		return nil
	}
	fmt.Printf("Pokémon in %s:\n", locationArea)
	for _, encounter := range result.PokemonEncounters {
		fmt.Println("-", encounter.Pokemon.Name)
	}
	return nil
}
