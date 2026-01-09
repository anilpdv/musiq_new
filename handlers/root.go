package handlers

import (
	"net/http"

	"musiq/models"

	"github.com/gin-gonic/gin"
)

// Root handles the root endpoint
func Root(c *gin.Context) {
	response := models.RootResponse{
		Status: http.StatusOK,
		Routes: models.Routes{
			SearchRoute:       "/api/search/:q",
			ListenRoute:       "/api/listen/:id/:name",
			WatchRoute:        "/api/watch/:id/:name",
			InfoRoute:         "/api/info/:id",
			RelatedRoute:      "/api/getvideo/:id",
			PlaylistRoute:     "/api/playlist/search/:q",
			PlaylistRouteByID: "/api/getplaylist/:id",
		},
	}

	c.JSON(http.StatusOK, response)
}
