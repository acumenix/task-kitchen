package api

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Error   string      `json:"error,omitempty"`
	Results interface{} `json:"results,omitempty"`
}

func SetupRouter(r *gin.RouterGroup, awsRegion, tableName string) {
	mgr := newKitchenManager(awsRegion, tableName)

	// Report endpoints
	r.GET("/:user", func(c *gin.Context) {
		fetchReportHandler(c, &mgr)
	})
	r.GET("/:user/:date", func(c *gin.Context) {
		getReportHandler(c, &mgr)
	})
	r.PUT("/:user/:date", func(c *gin.Context) {
		updateReportHandler(c, &mgr)
	})
	r.DELETE("/:user/:date", func(c *gin.Context) {
		deleteReportHandler(c, &mgr)
	})

	// Task Endpoint
	r.GET("/:user/:date/task", func(c *gin.Context) {
		getTasksHandler(c, &mgr)
	})
	r.POST("/:user/:date/task", func(c *gin.Context) {
		createTaskHandler(c, &mgr)
	})
	r.PUT("/:user/:date/task/:task_id", func(c *gin.Context) {
		updateTaskHandler(c, &mgr)
	})
	r.DELETE("/:user/:date/task/:task_id", func(c *gin.Context) {
		deleteTaskHandler(c, &mgr)
	})

	// Chore endpoints
	r.GET("/:user/:date/chore", func(c *gin.Context) {
		fetchChoresHandler(c, &mgr)
	})
	r.POST("/:user/:date/chore", func(c *gin.Context) {
		createChoreHandler(c, &mgr)
	})
	r.PUT("/:user/:date/chore/:chore_id", func(c *gin.Context) {
		updateChoreHandler(c, &mgr)
	})
	r.DELETE("/:user/:date/chore/:chore_id", func(c *gin.Context) {
		deleteChoreHandler(c, &mgr)
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
