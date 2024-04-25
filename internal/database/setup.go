package database

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/KhiemNguyen15/pokemon-card-api/internal/config"
	"github.com/KhiemNguyen15/pokemon-card-api/internal/models"
)

func PopulateCardDatabase(cfg config.PokemonAPIConfigurations) error {
	baseUrl := "https://api.pokemontcg.io/v2/cards?page="

	pageCount := cfg.PageCount
	for i := 1; i <= pageCount; i++ {
		url := baseUrl + fmt.Sprint(i)

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		request.Header.Add("X-Api-Key", cfg.APIKey)

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

			var cardSet models.Set
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

			var card models.Card
			card.Name = name
			card.Number = number
			card.Rarity = rarity
			card.ImageURL = imageUrl
			card.Value = cardValue
			card.SetName = cardSet.Name
			card.SetSeries = cardSet.Series

			if err := InsertCard(card); err != nil {
				return err
			}
		}
	}

	return nil
}

func PopulateSetDatabase(cfg config.PokemonAPIConfigurations) error {
	url := "https://api.pokemontcg.io/v2/sets"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Add("X-Api-Key", cfg.APIKey)

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
		var set models.Set
		set.Name = cardSet.(map[string]any)["name"].(string)
		set.Series = cardSet.(map[string]any)["series"].(string)
		set.Total = int(cardSet.(map[string]any)["total"].(float64))

		if err := InsertSet(set); err != nil {
			return err
		}
	}

	return nil
}
