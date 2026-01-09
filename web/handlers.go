package web

import (
	"musiq/services"
	"musiq/web/templates/components"
	"musiq/web/templates/pages"

	"github.com/gin-gonic/gin"
)

var youtubeService = services.NewYouTubeService()

// HomePage renders the main page
func HomePage(c *gin.Context) {
	pages.Home().Render(c.Request.Context(), c.Writer)
}

// SearchResultsView returns search results as HTML partial
func SearchResultsView(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		components.VideoGrid(nil).Render(c.Request.Context(), c.Writer)
		return
	}

	results, err := youtubeService.SearchVideos(query)
	if err != nil {
		components.VideoGrid(nil).Render(c.Request.Context(), c.Writer)
		return
	}

	components.VideoGrid(results).Render(c.Request.Context(), c.Writer)
}

// PlayerView returns the audio/video player as HTML partial
func PlayerView(c *gin.Context) {
	videoID := c.Param("id")
	playerType := c.DefaultQuery("type", "audio")

	// Get video info for title and author
	info, err := youtubeService.GetVideoInfo(videoID)
	title := "Unknown"
	author := "Unknown"
	if err == nil {
		title = info.Title
		author = info.Author
	}

	components.Player(videoID, title, author, playerType).Render(c.Request.Context(), c.Writer)
}

// PlaylistSearchView returns playlist search results as HTML partial
func PlaylistSearchView(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		components.PlaylistGrid(nil).Render(c.Request.Context(), c.Writer)
		return
	}

	results, err := youtubeService.SearchPlaylists(query)
	if err != nil {
		components.PlaylistGrid(nil).Render(c.Request.Context(), c.Writer)
		return
	}

	components.PlaylistGrid(results).Render(c.Request.Context(), c.Writer)
}

// PlaylistVideosView returns videos in a playlist as HTML partial
func PlaylistVideosView(c *gin.Context) {
	playlistID := c.Param("id")

	videos, err := youtubeService.GetPlaylistVideos(playlistID)
	if err != nil {
		components.PlaylistVideos(nil).Render(c.Request.Context(), c.Writer)
		return
	}

	components.PlaylistVideos(videos).Render(c.Request.Context(), c.Writer)
}
