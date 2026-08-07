package http

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(handlers *Handlers, jwtSecret []byte) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	api := r.Group("/api/v1")
	api.Use(AuthMiddleware(jwtSecret))
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("", handlers.GetTodayTasks)
			tasks.GET("/:task_id", handlers.GetTask)
			tasks.POST("/:task_id/complete", handlers.CompleteTask)
		}
	}
	return r
}
