package main

import (
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/m-mizutani/task-kitchen/api"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	logger.SetLevel(logrus.DebugLevel)
	api.Logger = logger

	if len(os.Args) != 4 {
		logger.Fatal("syntax error) server [region] [table_name] [static_dir]")
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	api.SetupRouter(v1, os.Args[1], os.Args[2])
	r.Use(static.Serve("/", static.LocalFile(os.Args[3], false)))

	r.Run()
}
