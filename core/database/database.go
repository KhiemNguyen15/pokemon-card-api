package database

import (
	"database/sql"

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
