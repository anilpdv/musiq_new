package handlers

import (
	"log"
	"net/http"

	"musiq/models"
	"musiq/services"

	"github.com/gin-gonic/gin"
)

// Info handles video info requests
func Info(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing video ID",
		})
		return
	}

	// Extract video ID from URL if necessary
	videoID = services.ExtractVideoID(videoID)

	info, err := youtubeService.GetVideoInfo(videoID)
	if err != nil {
		log.Printf("Failed to get video info for %s: %v", videoID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get video info",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, info)
}
