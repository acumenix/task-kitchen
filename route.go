package kitchen

import (
	"time"

	"github.com/gin-gonic/gin"
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

func BuildRouter(awsRegion, tableName string) *gin.Engine {
	mgr := newKitchenManager(awsRegion, tableName)

	r := gin.Default()
	r.GET("/v1/:user/:date/", func(c *gin.Context) {
		user := getParam(c.Params, "user")
		date := getParam(c.Params, "date")
		ts, err := time.Parse("2006-01-02", date)
		if err != nil {
			Logger.WithError(err).Error("Format error")
			c.JSON(400, Response{"Invalid date format, should be like 2006-01-02", nil})
			return
		}

		tasks, err := mgr.FetchTasks(user, ts)
		if err != nil {
			Logger.WithError(err).Error("Fail to fetch tasks")
			c.JSON(500, Response{"Internal server error", nil})
			return
		}

		c.JSON(200, Response{"", tasks})
	})

	// Task Endpoint
	r.GET("/v1/:user/:date/task", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "get tasks",
		})
	})

	r.POST("/v1/:user/:date/task", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "create a task",
		})
	})

	r.PUT("/v1/:user/:date/task/:task_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.DELETE("/v1/:user/:date/task/:task_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Pomodoro Endpoint
	r.GET("/v1/:user/:date/pomodoro/:task_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/v1/:user/:date/pomodoro/:task_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.PUT("/v1/:user/:date/pomodoro/:task_id/:pomodoro_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.DELETE("/v1/:user/:date/pomodoro/:task_id/:pomodoro_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "delete pomodoro",
		})
	})

	return r
}
