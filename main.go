package main

import (
	"pokemon-api/configs"
	"pokemon-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	configs.ConnectDB()
	routes.UserRoute(r)
	r.Run()
}
