package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.RouterGroup, awsRegion, tableName string) {
	mgr := newKitchenManager(awsRegion, tableName)

	// Report endpoints
	r.GET("/:user", func(c *gin.Context) {
		handle(fetchReportHandler, c, &mgr)
	})
	r.GET("/:user/:date", func(c *gin.Context) {
		handle(getReportHandler, c, &mgr)
	})
	r.PUT("/:user/:date", func(c *gin.Context) {
		handle(updateReportHandler, c, &mgr)
	})
	r.DELETE("/:user/:date", func(c *gin.Context) {
		handle(deleteReportHandler, c, &mgr)
	})

	// Task Endpoint
	r.GET("/:user/:date/task", func(c *gin.Context) {
		handle(getTasksHandler, c, &mgr)
	})
	r.POST("/:user/:date/task", func(c *gin.Context) {
		handle(createTaskHandler, c, &mgr)
	})
	r.PUT("/:user/:date/task/:task_id", func(c *gin.Context) {
		handle(updateTaskHandler, c, &mgr)
	})
	r.DELETE("/:user/:date/task/:task_id", func(c *gin.Context) {
		handle(deleteTaskHandler, c, &mgr)
	})

	// Chore endpoints
	r.GET("/:user/:date/chore", func(c *gin.Context) {
		handle(fetchChoresHandler, c, &mgr)
	})
	r.POST("/:user/:date/chore", func(c *gin.Context) {
		handle(createChoreHandler, c, &mgr)
	})
	r.PUT("/:user/:date/chore/:chore_id", func(c *gin.Context) {
		handle(updateChoreHandler, c, &mgr)
	})
	r.DELETE("/:user/:date/chore/:chore_id", func(c *gin.Context) {
		handle(deleteChoreHandler, c, &mgr)
	})

	// Pomodoro Endpoint
	r.GET("/:user/:date/pomodoro/:task_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/:user/:date/pomodoro/:task_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.PUT("/:user/:date/pomodoro/:task_id/:pomodoro_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.DELETE("/:user/:date/pomodoro/:task_id/:pomodoro_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "delete pomodoro",
		})
	})
}
