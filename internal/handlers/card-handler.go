package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/KhiemNguyen15/pokemon-card-api/internal/database"
)

func GetCards(c *gin.Context) {
	rarity := c.Query("rarity")
	if rarity != "" {
		cards, err := database.GetCardsByRarity(rarity)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "cards not found"})
			return
		}

		c.JSON(http.StatusOK, cards)
		return
	}

	cards, err := database.GetCards()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cards not found"})
		return
	}

	c.JSON(http.StatusOK, cards)
}

func GetCardByID(c *gin.Context) {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}

	card, err := database.GetCard(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "card not found"})
		return
	}

	c.JSON(http.StatusOK, card)
}
