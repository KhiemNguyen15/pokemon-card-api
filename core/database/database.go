package database

import (
	"database/sql"

	"core/card"
	"github.com/go-sql-driver/mysql"
	"utils/config"
)

func LoadDatabase(dbConfig config.DatabaseConfigurations) (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 dbConfig.DBUser,
		Passwd:               dbConfig.DBPassword,
		Net:                  "tcp",
		Addr:                 dbConfig.DBHost,
		DBName:               dbConfig.DBName,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetCard(db *sql.DB, id int) (card.Card, error) {
	var card card.Card

	row := db.QueryRow(
		"SELECT (id, name, number, rarity, value, image_url, card_set) FROM pokemon_cards WHERE id = ?",
		id,
	)
	err := row.Scan(
		&card.ID,
		&card.Name,
		&card.Number,
		&card.Rarity,
		&card.Value,
		&card.ImageURL,
		&card.Set.Name,
	)
	if err != nil {
		return card, err
	}

	return card, nil
}

func InsertCard(db *sql.DB, card card.Card) error {
	_, err := db.Exec(
		"INSERT INTO pokemon_cards (name, number, rarity, value, image_url, card_set)"+
			"VALUES (?, ?, ?, ?, ?, ?)",
		card.Name,
		card.Number,
		card.Rarity,
		card.Value,
		card.ImageURL,
		card.Set.Name,
	)
	if err != nil {
		return err
	}

	return nil
}

func InsertSet(db *sql.DB, set card.Set) error {
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
