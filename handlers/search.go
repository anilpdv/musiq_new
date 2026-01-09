package handlers

import (
	"log"
	"net/http"

	"musiq/models"
	"musiq/services"

	"github.com/gin-gonic/gin"
)

var youtubeService = services.NewYouTubeService()

// Search handles video search requests
func Search(c *gin.Context) {
	query := c.Param("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing search query",
		})
		return
	}

	videos, err := youtubeService.SearchVideos(query)
	if err != nil {
		log.Printf("Search error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Search failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, videos)
}
