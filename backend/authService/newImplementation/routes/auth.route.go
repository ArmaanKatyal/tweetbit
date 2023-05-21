package routes

import (
	auth "github.com/ArmaanKatyal/tweetbit/backend/authService/controllers"
	"github.com/ArmaanKatyal/tweetbit/backend/authService/middlewares"
	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	rg *gin.RouterGroup
}

func NewAuthRoute(router *gin.RouterGroup) *AuthRoute {
	return &AuthRoute{
		rg: router,
	}
}

func (a *AuthRoute) HandleRoutes() {
	a.login()
	a.logout()
	a.refresh()
	a.register()
}

func (a *AuthRoute) login() {
	a.rg.POST("/login", auth.Login)
}

func (a *AuthRoute) logout() {
	a.rg.POST("/logout", auth.Logout)
}

func (a *AuthRoute) refresh() {
	rg := a.rg.Group("/refresh", middlewares.VerifyToken())
	rg.GET("/", auth.Refresh)
}

func (a *AuthRoute) register() {
	a.rg.POST("/register", auth.Register)
}
