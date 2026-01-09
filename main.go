package main

import (
	"log"
	"os"

	"musiq/handlers"
	"musiq/middleware"
	"musiq/web"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Apply CORS middleware
	r.Use(middleware.CORS())

	// Serve static files
	r.Static("/static", "./web/static")

	// UI routes (templ-rendered HTML)
	r.GET("/", web.HomePage)
	ui := r.Group("/ui")
	{
		ui.GET("/search", web.SearchResultsView)
		ui.GET("/play/:id", web.PlayerView)
		ui.GET("/playlists", web.PlaylistSearchView)
		ui.GET("/playlist/:id", web.PlaylistVideosView)
	}

	// API routes (JSON)
	api := r.Group("/api")
	{
		// Search
		api.GET("/search/:q", handlers.Search)

		// Audio/Video streaming
		api.GET("/listen/:id/:name", handlers.Listen)
		api.GET("/watch/:id/:name", handlers.Watch)

		// Video info
		api.GET("/info/:id", handlers.Info)

		// Related videos
		api.GET("/getvideo/:id", handlers.GetVideo)
		api.GET("/related/:id", handlers.Related)

		// Playlists
		api.GET("/playlist/search/:q", handlers.PlaylistSearch)
		api.GET("/getplaylist/:id", handlers.GetPlaylist)
	}

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
