package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/KhiemNguyen15/pokemon-card-api/internal/config"
	"github.com/KhiemNguyen15/pokemon-card-api/internal/database"
	"github.com/KhiemNguyen15/pokemon-card-api/internal/handlers"
)

var (
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	WarningLogger *log.Logger
	DebugLogger   *log.Logger
)

func init() {
	logFilePath := "./data/logs.txt"

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}

	// gin.DisableConsoleColor()
	// gin.DefaultWriter = logFile
	// gin.DefaultErrorWriter = logFile

	InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(logFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	InfoLogger.Println("Process starting...")

	configPath := "./data/config.yml"

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		ErrorLogger.Fatalln("Error loading config file:", err)
	}

	err = database.ConnectDatabase(cfg.Database)
	if err != nil {
		ErrorLogger.Fatalln("Error connecting to database:", err)
	}

	// Populate the databases (only do this once)
	// if err := database.PopulateSetDatabase(cfg.PokemonAPI); err != nil {
	// 	panic(err)
	// }
	//
	// if err := database.PopulateCardDatabase(cfg.PokemonAPI); err != nil {
	// 	panic(err)
	// }

	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/cards", handlers.GetCards)
	router.GET("/cards/:id", handlers.GetCardByID)
	router.GET("/sets", handlers.GetSets)
	router.GET("/sets/:name", handlers.GetSetByName)

	router.Run(":8080")
}
