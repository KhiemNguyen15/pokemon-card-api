package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/KhiemNguyen15/pokemon-card-api/internal/database"
)

func GetSets(c *gin.Context) {
	sets, err := database.GetSets()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sets not found"})
		return
	}

	c.JSON(http.StatusOK, sets)
}

func GetSetByName(c *gin.Context) {
	name := c.Param("name")

	set, err := database.GetSet(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "set not found"})
		return
	}

	c.JSON(http.StatusOK, set)
}
