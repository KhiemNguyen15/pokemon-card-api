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

	row := db.QueryRow("SELECT (id, name, value, image_url) FROM pokemon_cards WHERE id = ?", id)
	if err := row.Scan(&card.ID, &card.Name, &card.Value, &card.ImageURL); err != nil {
		return card, err
	}

	return card, nil
}

func InsertCard(db *sql.DB, card card.Card) error {
	_, err := db.Exec(
		"INSERT INTO pokemon_cards (name, value, image_url) VALUES (?, ?, ?)",
		card.Name,
		card.Value,
		card.ImageURL,
	)
	if err != nil {
		return err
	}

	return nil
}
