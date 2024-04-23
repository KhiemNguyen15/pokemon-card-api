package database

import (
	"github.com/KhiemNguyen15/PokemonCardTrader/internal/models"
)

func GetCard(id int) (models.Card, error) {
	card := models.Card{}

	err := DB.Get(&card, "SELECT * FROM pokemon_cards WHERE id = ?", id)
	if err != nil {
		return card, err
	}

	return card, nil
}

func GetCards() ([]models.Card, error) {
	cards := []models.Card{}

	err := DB.Select(&cards, "SELECT * FROM pokemon_cards")
	if err != nil {
		return cards, err
	}

	return cards, nil
}

func GetCardsByRarity(rarity string) ([]models.Card, error) {
	cards := []models.Card{}

	err := DB.Select(&cards, "SELECT * FROM pokemon_cards WHERE rarity = ?", rarity)
	if err != nil {
		return cards, err
	}

	return cards, nil
}

func InsertCard(card models.Card) error {
	_, err := DB.Exec(
		"INSERT INTO pokemon_cards "+
			"(name, number, rarity, value, image_url, set_name, set_series) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?)",
		card.Name,
		card.Number,
		card.Rarity,
		card.Value,
		card.ImageURL,
		card.SetName,
		card.SetSeries,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetSet(name string) (models.Set, error) {
	set := models.Set{}

	err := DB.Get(&set, "SELECT * FROM card_sets WHERE name = ?", name)
	if err != nil {
		return set, err
	}

	return set, nil
}

func GetSets() ([]models.Set, error) {
	sets := []models.Set{}

	err := DB.Select(&sets, "SELECT * FROM card_sets")
	if err != nil {
		return sets, err
	}

	return sets, nil
}

func InsertSet(set models.Set) error {
	_, err := DB.Exec(
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
