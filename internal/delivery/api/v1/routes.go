package v1

import "github.com/gin-gonic/gin"

func (h *Handler) Routes() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	router.Use(gin.Recovery())
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.Login)
			auth.PUT("/refresh", h.Refresh)
			auth.GET("/sessions", h.GetAllSessions)
		}
	}
	return router
}
