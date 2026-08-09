package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(handlers *Handlers, jwtSecret []byte) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	api := r.Group("/api/v1")
	api.Use(AuthMiddleware(jwtSecret))
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("", handlers.GetTodayTasks)
			tasks.GET("/:task_id", handlers.GetTask)
			tasks.PUT("/:task_id/complete", handlers.CompleteTask)
		}
	}
	internal := r.Group("/internal")
	{
		internal.GET("/tasks/:user_id", handlers.GetTodayTasksForInternal)
	}
	return r
}
