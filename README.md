# Pokémon CLI Application

A command-line interface (CLI) application that allows users to explore the Pokémon world, catch Pokémon, and manage their Pokédex.

## Features

- **Map Exploration**: View different location areas in the Pokémon world
- **Area Exploration**: Discover Pokémon in specific location areas
- **Pokémon Catching**: Attempt to catch Pokémon with a chance-based system
- **Pokédex Management**: View and inspect caught Pokémon
- **Interactive CLI**: User-friendly command-line interface

## Prerequisites

- Go 1.16 or higher
- Internet connection (for API calls)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd Pokedex
```

2. Build the application:
```bash
go build
```

## Usage

Run the application:
```bash
./Pokedex
```

### Available Commands

- `help`: Display help message
- `exit`: Exit the application
- `map`: Display 20 location areas in the Pokémon world
- `explore <location-area>`: List all Pokémon in a specific location area
- `catch <pokemon>`: Attempt to catch a Pokémon
- `inspect <pokemon>`: View details of a caught Pokémon
- `pokedex`: List all caught Pokémon

### Examples

1. View location areas:
```
Pokedex > map
```

2. Explore a specific area:
```
Pokedex > explore pallet-town
```

3. Catch a Pokémon:
```
Pokedex > catch pikachu
```

4. Inspect a caught Pokémon:
```
Pokedex > inspect pikachu
```

## How It Works

- The application uses the [PokéAPI](https://pokeapi.co/) to fetch Pokémon and location data
- Catching Pokémon is based on a chance system that considers the Pokémon's base experience
- The Pokédex stores information about caught Pokémon locally during the session

## Contributing

Feel free to submit issues and enhancement requests!

## License

This project is open source and available under the MIT License.

