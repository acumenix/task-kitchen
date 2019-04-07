package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"

	"github.com/gin-gonic/gin"
	"github.com/m-mizutani/task-kitchen/api"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	api.Logger = logger

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	v1 := r.Group("/v1")
	api.SetupRouter(v1, os.Getenv("AWS_REGION"), os.Getenv("TABLE_NAME"))
	ginLambda := ginadapter.New(r)

	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return ginLambda.Proxy(req)
	})
}
