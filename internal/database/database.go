package database

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/KhiemNguyen15/PokemonCardTrader/internal/card"
	"github.com/KhiemNguyen15/PokemonCardTrader/internal/config"
)

func LoadDatabase(dbConfig config.DatabaseConfigurations) (*sqlx.DB, error) {
	cfg := mysql.Config{
		User:                 dbConfig.DBUser,
		Passwd:               dbConfig.DBPassword,
		Net:                  "tcp",
		Addr:                 dbConfig.DBHost,
		DBName:               dbConfig.DBName,
		AllowNativePasswords: true,
	}

	db, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetCard(db *sqlx.DB, id int) (card.Card, error) {
	c := card.Card{}

	err := db.Get(&c, "SELECT * FROM pokemon_cards WHERE id = ?", id)
	if err != nil {
		return c, err
	}

	return c, nil
}

func GetCards(db *sqlx.DB) ([]card.Card, error) {
	cards := []card.Card{}

	err := db.Select(&cards, "SELECT * FROM pokemon_cards")
	if err != nil {
		return cards, err
	}

	return cards, nil
}

func PopulateCardDatabase(db *sqlx.DB, cfg config.PokemonAPIConfigurations) error {
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

			if err := insertCard(db, card); err != nil {
				return err
			}
		}
	}

	return nil
}

func PopulateSetDatabase(db *sqlx.DB, cfg config.PokemonAPIConfigurations) error {
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
		var set card.Set
		set.Name = cardSet.(map[string]any)["name"].(string)
		set.Series = cardSet.(map[string]any)["series"].(string)
		set.Total = int(cardSet.(map[string]any)["total"].(float64))

		if err := insertSet(db, set); err != nil {
			return err
		}
	}

	return nil
}

func insertCard(db *sqlx.DB, c card.Card) error {
	_, err := db.Exec(
		"INSERT INTO pokemon_cards "+
			"(name, number, rarity, value, image_url, set_name, set_series) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?)",
		c.Name,
		c.Number,
		c.Rarity,
		c.Value,
		c.ImageURL,
		c.Set.Name,
		c.Set.Series,
	)
	if err != nil {
		return err
	}

	return nil
}

func insertSet(db *sqlx.DB, set card.Set) error {
	_, err := db.Exec(
		"INSERT INTO card_sets (name, series, total) VALUES (?, ?, ?)",
		set.Name,
		set.Series,
		set.Total,
	)
	if err != nil {
		return err
	}

	return nil
}
