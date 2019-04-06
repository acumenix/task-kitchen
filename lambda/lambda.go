package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/gin-gonic/gin"
	kitchen "github.com/m-mizutani/task-kitchen"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	kitchen.Logger = logger

	r := gin.Default()
	v1 := r.Group("/v1")
	kitchen.SetupRouter(v1, os.Getenv("AWS_REGION"), os.Getenv("TABLE_NAME"))
	ginLambda := ginadapter.New(r)

	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return ginLambda.Proxy(req)
	})
}
