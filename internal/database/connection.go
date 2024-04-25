package database

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/KhiemNguyen15/pokemon-card-api/internal/config"
)

var DB *sqlx.DB

func ConnectDatabase(dbConfig config.DatabaseConfigurations) error {
	cfg := mysql.Config{
		User:                 dbConfig.DBUser,
		Passwd:               dbConfig.DBPassword,
		Net:                  "tcp",
		Addr:                 dbConfig.DBHost,
		DBName:               dbConfig.DBName,
		AllowNativePasswords: true,
	}

	database, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	if err := database.Ping(); err != nil {
		return err
	}

	DB = database

	return nil
}
