package api

import (
	"github.com/gin-gonic/gin"
	"shopify-integration/pkg/model"

	"net/http"
)

func configureRouter(api *api) *gin.Engine {
	router := gin.Default()

	router.Use(CORSMiddleware())

	public := router.Group("api/v1")

	public.POST("/login", api.Auth().Login)
	public.POST("/refresh", api.Auth().Refresh)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, model.ErrRecordNotFound)
	})

	return router
}
