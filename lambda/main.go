package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"

	kitchen "github.com/m-mizutani/task-kitchen"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	r := kitchen.BuildRouter(os.Getenv("AWS_REGION"), os.Getenv("TABLE_NAME"))
	ginLambda := ginadapter.New(r)

	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return ginLambda.Proxy(req)
	})
}
