package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/KhiemNguyen15/PokemonCardTrader/internal/config"
	"github.com/KhiemNguyen15/PokemonCardTrader/internal/database"
)

var (
	db *sqlx.DB

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

	gin.DisableConsoleColor()
	gin.DefaultWriter = logFile
	gin.DefaultErrorWriter = logFile

	InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(logFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	InfoLogger.Println("\nProcess starting...")

	configPath := "./data/config.yml"

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		ErrorLogger.Fatalln("Error loading config file:", err)
	}

	db, err = database.LoadDatabase(cfg.Database)
	if err != nil {
		ErrorLogger.Fatalln("Error connecting to database:", err)
	}

	// Populate the databases (only do this once)
	// if err := database.PopulateSetDatabase(db, config.PokemonAPI); err != nil {
	// 	panic(err)
	// }
	//
	// if err := database.PopulateCardDatabase(db, config.PokemonAPI); err != nil {
	// 	panic(err)
	// }

	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/cards", getCards)
	router.GET("/cards/:id", getCardByID)
	router.GET("/ping", ping)

	router.Run(":8080")
}

func getCards(c *gin.Context) {
	cards, err := database.GetCards(db)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cards not found"})
		ErrorLogger.Println("Error retrieving cards from database:", err)
		return
	}

	c.JSON(http.StatusOK, cards)
}

func getCardByID(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}

	card, err := database.GetCard(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}

	c.JSON(http.StatusOK, card)
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
