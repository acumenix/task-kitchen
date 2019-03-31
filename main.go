package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

var ginLambda *ginadapter.GinLambda

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.Proxy(req)
}

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	r := gin.Default()
	r.GET("/v1/:user/:date/", func(c *gin.Context) {
		logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Task Endpoint
	r.GET("/v1/:user/:date/task", func(c *gin.Context) {
		logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/v1/:user/:date/task", func(c *gin.Context) {
		logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.PUT("/v1/:user/:date/task/:task_id", func(c *gin.Context) {
		logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.DELETE("/v1/:user/:date/task/:task_id", func(c *gin.Context) {
		logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Pomodoro Endpoint
	r.GET("/v1/:user/:date/pomodoros/:task_id", func(c *gin.Context) {
		logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/v1/:user/:date/pomodoros/:task_id", func(c *gin.Context) {
		logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.PUT("/v1/:user/:date/pomodoros/:task_id/:pomodoro_id", func(c *gin.Context) {
		logger.WithField("param", c.Params).Info("Request")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	ginLambda = ginadapter.New(r)

	lambda.Start(handler)
}
