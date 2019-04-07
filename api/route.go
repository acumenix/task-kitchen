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

	r.GET("/:user/:date/", func(c *gin.Context) {
		user, ts, err := getSpace(c.Params)
		if err != nil {
			c.JSON(400, Response{err.Error(), nil})
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
	r.GET("/:user/:date/task", func(c *gin.Context) {
		user, ts, err := getSpace(c.Params)
		if err != nil {
			c.JSON(400, Response{err.Error(), nil})
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

	r.POST("/:user/:date/task", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		user, ts, err := getSpace(c.Params)
		if err != nil {
			c.JSON(400, Response{err.Error(), nil})
			return
		}

		task, err := mgr.NewTask(user, ts)
		if err != nil {
			Logger.WithError(err).Error("Fail to create a task")
			c.JSON(500, Response{"Internal server error", nil})
			return
		}

		c.JSON(200, Response{"ok", task})
	})

	r.PUT("/:user/:date/task/:task_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.DELETE("/:user/:date/task/:task_id", func(c *gin.Context) {
		Logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
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
