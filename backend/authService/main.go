package main

import (
	"github.com/ArmaanKatyal/tweetbit/backend/authService/routes"
	"github.com/ArmaanKatyal/tweetbit/backend/authService/services"
)

func main() {
	ginService := services.NewGinService()
	rg := ginService.Engine.Group("/api")
	authRoute := routes.NewAuthRoute(rg)
	authRoute.HandleRoutes()
	ginService.Run()
}
