package main

import (
	"pokemon-api/configs"
	"pokemon-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	configs.ConnectDB()
	routes.UserRoute(r)
	configCors := cors.DefaultConfig()
	configCors.AllowAllOrigins = true
	r.Run()
}
