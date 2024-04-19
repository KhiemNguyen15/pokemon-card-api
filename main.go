package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/KhiemNguyen15/PokemonCardTrader/internal/config"
	"github.com/KhiemNguyen15/PokemonCardTrader/internal/database"
)

var db *sqlx.DB

func main() {
	configPath := "./data/config.yml"

	config, err := config.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Errorf("Error loading config file: %v\n", err))
	}

	db, err = database.LoadDatabase(config.Database)
	if err != nil {
		panic(fmt.Errorf("Error connecting to database: %v\n", err))
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

	fmt.Println("Process ended gracefully.")
}

func getCards(c *gin.Context) {
	cards, err := database.GetCards(db)
	if err != nil {
		panic(fmt.Errorf("Error retrieving cards from database: %v\n", err))
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
	c.JSON(http.StatusOK, gin.H{"error": "pong"})
}
