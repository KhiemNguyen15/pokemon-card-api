package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/KhiemNguyen15/PokemonCardTrader/internal/card"
	"github.com/KhiemNguyen15/PokemonCardTrader/internal/config"
	"github.com/KhiemNguyen15/PokemonCardTrader/internal/database"
)

func populateCardDatabase(db *sql.DB, config config.PokemonAPIConfigurations) error {
	baseUrl := "https://api.pokemontcg.io/v2/cards?page="

	pageCount := config.PageCount
	for i := 1; i <= pageCount; i++ {
		url := baseUrl + fmt.Sprint(i)

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		request.Header.Add("X-Api-Key", config.APIKey)

		client := http.Client{}
		response, err := client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("Unexpected status code: %v\n", response.StatusCode)
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var result map[string]any
		json.Unmarshal(body, &result)

		data := result["data"].([]any)
		for _, pokemonCard := range data {
			marketData := pokemonCard.(map[string]any)["tcgplayer"]
			if marketData == nil {
				continue
			}
			priceData := marketData.(map[string]any)["prices"]
			if priceData == nil {
				continue
			}

			var cardValue float64 = 0
			for _, value := range priceData.(map[string]any) {
				if value.(map[string]any)["market"] != nil {
					cardValue = value.(map[string]any)["market"].(float64)
					break
				} else if value.(map[string]any)["high"] != nil {
					cardValue = value.(map[string]any)["high"].(float64)
					break
				}
			}
			if cardValue == 0 {
				continue
			}

			var cardSet card.Set
			set := pokemonCard.(map[string]any)["set"].(map[string]interface{})
			cardSet.Name = set["name"].(string)
			cardSet.Series = set["series"].(string)
			cardSet.Total = int(set["total"].(float64))

			rarityOptional := pokemonCard.(map[string]any)["rarity"]
			var rarity string
			if rarityOptional == nil {
				rarity = "N/A"
			} else {
				rarity = rarityOptional.(string)
			}

			name := pokemonCard.(map[string]any)["name"].(string)
			number := pokemonCard.(map[string]any)["number"].(string)
			imageUrl := pokemonCard.(map[string]any)["images"].(map[string]any)["small"].(string)

			var card card.Card
			card.Name = name
			card.Number = number
			card.Rarity = rarity
			card.ImageURL = imageUrl
			card.Value = cardValue
			card.Set = cardSet

			if err := database.InsertCard(db, card); err != nil {
				return err
			}
		}
	}

	return nil
}

func populateSetDatabase(db *sql.DB, config config.PokemonAPIConfigurations) error {
	url := "https://api.pokemontcg.io/v2/sets"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Add("X-Api-Key", config.APIKey)

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: %v\n", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	data := result["data"].([]any)
	for _, cardSet := range data {
		var set card.Set
		set.Name = cardSet.(map[string]any)["name"].(string)
		set.Series = cardSet.(map[string]any)["series"].(string)
		set.Total = int(cardSet.(map[string]any)["total"].(float64))

		if err := database.InsertSet(db, set); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	configPath := "./data/config.yml"

	config, err := config.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Errorf("Error loading config file: %v\n", err))
	}

	db, err := database.LoadDatabase(config.Database)
	if err != nil {
		panic(fmt.Errorf("Error connecting to database: %v\n", err))
	}

	fmt.Println("Idle database connections: ", db.Stats().Idle)

	bot, err := discordgo.New("Bot " + config.DiscordAPI.APIKey)
	if err != nil {
		panic(fmt.Errorf("Error creating Discord session: %v\n", err))
	}

	if err := bot.Open(); err != nil {
		panic(fmt.Errorf("Error while opening Discord session: %v\n", err))
	}

	fmt.Println("Bot successfully started.")

	defer bot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, os.Interrupt)
	<-stop

	fmt.Println("Process ended gracefully.")

	// Populate the databases (only do this one)
	// if err := populateSetDatabase(db, config.PokemonAPI); err != nil {
	// 	panic(err)
	// }
	//
	// if err := populateCardDatabase(db, config.PokemonAPI); err != nil {
	// 	panic(err)
	// }
}
