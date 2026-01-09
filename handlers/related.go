package handlers

import (
	"log"
	"net/http"

	"musiq/models"
	"musiq/services"

	"github.com/gin-gonic/gin"
)

// GetVideo handles related videos request (using search as fallback)
func GetVideo(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing video ID",
		})
		return
	}

	// Extract video ID from URL if necessary
	videoID = services.ExtractVideoID(videoID)

	// Get video info first to use title for related search
	info, err := youtubeService.GetVideoInfo(videoID)
	if err != nil {
		log.Printf("Failed to get video info for %s: %v", videoID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Something went wrong",
			Message: err.Error(),
		})
		return
	}

	// Search for related videos using the video title
	relatedVideos, err := youtubeService.SearchVideos(info.Title)
	if err != nil {
		log.Printf("Failed to get related videos for %s: %v", videoID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Something went wrong",
			Message: err.Error(),
		})
		return
	}

	// Filter out the original video from results
	filteredResults := make([]models.VideoResult, 0, len(relatedVideos))
	for _, v := range relatedVideos {
		if v.ID != videoID {
			filteredResults = append(filteredResults, v)
		}
	}

	c.JSON(http.StatusOK, filteredResults)
}

// Related handles video details + related videos request
func Related(c *gin.Context) {
	videoID := c.Param("id")
	if videoID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing video ID",
		})
		return
	}

	// Extract video ID from URL if necessary
	videoID = services.ExtractVideoID(videoID)

	// Get video info
	info, err := youtubeService.GetVideoInfo(videoID)
	if err != nil {
		log.Printf("Failed to get video info for %s: %v", videoID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Something went wrong",
			Message: err.Error(),
		})
		return
	}

	// Search for related videos using the video title
	relatedVideos, err := youtubeService.SearchVideos(info.Title)
	if err != nil {
		log.Printf("Failed to get related videos for %s: %v", videoID, err)
		// Return video details even if related fails
		response := models.RelatedResponse{
			VideoDetails: *info,
			RelatedSongs: []models.VideoResult{},
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// Filter out the original video
	filteredResults := make([]models.VideoResult, 0, len(relatedVideos))
	for _, v := range relatedVideos {
		if v.ID != videoID {
			filteredResults = append(filteredResults, v)
		}
	}

	response := models.RelatedResponse{
		VideoDetails: *info,
		RelatedSongs: filteredResults,
	}

	c.JSON(http.StatusOK, response)
}
