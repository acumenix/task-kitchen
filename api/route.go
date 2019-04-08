package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Response struct {
	Error   string      `json:"error,omitempty"`
	Results interface{} `json:"results,omitempty"`
}

func getParam(params gin.Params, key string) string {
	for _, p := range params {
		if p.Key == key {
			return p.Value
		}
	}

	return ""
}

func getSpace(params gin.Params) (user string, ts time.Time, err error) {
	if user = getParam(params, "user"); user == "" {
		err = errors.New("user parameter is empty")
		return
	}

	date := getParam(params, "date")
	if ts, err = time.Parse("2006-01-02", date); err != nil {
		Logger.WithError(err).Error("Format error")
		err = errors.New("Invalid date format, should be like 2006-01-02")
	}

	return
}

func SetupRouter(r *gin.RouterGroup, awsRegion, tableName string) {
	mgr := newKitchenManager(awsRegion, tableName)

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
