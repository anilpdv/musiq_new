package handlers

import (
	"log"
	"net/http"

	"musiq/models"

	"github.com/gin-gonic/gin"
)

// PlaylistSearch handles playlist search requests
func PlaylistSearch(c *gin.Context) {
	query := c.Param("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing search query",
		})
		return
	}

	playlists, err := youtubeService.SearchPlaylists(query)
	if err != nil {
		log.Printf("Playlist search error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Search failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, playlists)
}

// GetPlaylist handles get playlist videos request
func GetPlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	if playlistID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing playlist ID",
		})
		return
	}

	videos, err := youtubeService.GetPlaylistVideos(playlistID)
	if err != nil {
		log.Printf("Failed to get playlist %s: %v", playlistID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get playlist",
			Message: err.Error(),
		})
		return
	}

	if len(videos) == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "No items found",
		})
		return
	}

	c.JSON(http.StatusOK, videos)
}
