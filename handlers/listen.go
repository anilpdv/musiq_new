package handlers

import (
	"log"
	"net/http"

	"musiq/models"
	"musiq/services"

	"github.com/gin-gonic/gin"
)

var ffmpegService = services.NewFFmpegService()

// Listen handles MP3 audio streaming
func Listen(c *gin.Context) {
	videoID := c.Param("id")
	filename := c.Param("name")

	if videoID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing video ID",
		})
		return
	}

	// Extract video ID from URL if necessary
	videoID = services.ExtractVideoID(videoID)

	// Get audio stream from YouTube
	audioStream, _, err := youtubeService.GetAudioStream(videoID)
	if err != nil {
		log.Printf("Failed to get audio stream for %s: %v", videoID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get audio stream",
			Message: err.Error(),
		})
		return
	}
	defer audioStream.Close()

	// Set response headers
	c.Header("Content-Type", "audio/mpeg")
	c.Header("Cache-Control", "public, max-age=3600")

	// Set download header if download=true query param
	if c.Query("download") == "true" {
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	} else {
		c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	}

	// Convert to MP3 and stream to response
	if err := ffmpegService.ConvertToMP3(audioStream, c.Writer); err != nil {
		log.Printf("MP3 conversion error for %s: %v", videoID, err)
		// Only send error if headers haven't been sent
		if !c.Writer.Written() {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Conversion failed",
				Message: err.Error(),
			})
		}
		return
	}
}
