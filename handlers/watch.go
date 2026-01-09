package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"musiq/models"
	"musiq/services"

	"github.com/gin-gonic/gin"
)

// flushWriter wraps an io.Writer and flushes after each write
type flushWriter struct {
	w       io.Writer
	flusher http.Flusher
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.flusher != nil {
		fw.flusher.Flush()
	}
	return
}

// Watch handles MP4 video streaming with progressive download
func Watch(c *gin.Context) {
	videoID := c.Param("id")

	if videoID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Missing video ID",
		})
		return
	}

	// Extract video ID from URL if necessary
	videoID = services.ExtractVideoID(videoID)

	// Try to get a combined video+audio stream first (instant playback)
	stream, mimeType, size, err := youtubeService.GetCombinedStream(videoID)
	if err == nil {
		defer stream.Close()
		log.Printf("Using combined stream for %s (size: %d, type: %s)", videoID, size, mimeType)

		// Set headers for direct streaming
		c.Header("Content-Type", "video/mp4")
		if size > 0 {
			c.Header("Content-Length", fmt.Sprintf("%d", size))
		}
		c.Header("Accept-Ranges", "bytes")
		c.Header("Cache-Control", "no-cache")

		// Set download header if download=true query param
		filename := c.Param("name")
		if c.Query("download") == "true" {
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
		}

		// Stream directly to client
		c.Writer.WriteHeader(http.StatusOK)
		if _, err := io.Copy(c.Writer, stream); err != nil {
			log.Printf("Stream copy error for %s: %v", videoID, err)
		}
		return
	}

	// Fallback to muxing separate streams if no combined format available
	log.Printf("No combined stream for %s, falling back to mux: %v", videoID, err)

	videoStream, audioStream, _, err := youtubeService.GetVideoAndAudioStreams(videoID)
	if err != nil {
		log.Printf("Failed to get video streams for %s: %v", videoID, err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get video streams",
			Message: err.Error(),
		})
		return
	}
	defer videoStream.Close()
	defer audioStream.Close()

	// Set response headers for streaming
	c.Header("Content-Type", "video/mp4")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Flush()

	// Create a flushing writer to ensure data is sent immediately
	fw := &flushWriter{
		w:       c.Writer,
		flusher: c.Writer,
	}

	// Use pipe-based mux
	if err := ffmpegService.MuxVideoAudio(videoStream, audioStream, fw); err != nil {
		log.Printf("Video streaming error for %s: %v", videoID, err)
		return
	}
}
