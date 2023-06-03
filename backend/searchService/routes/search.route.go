package routes

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	rg *gin.RouterGroup
}

func NewAuthRoute(rg *gin.RouterGroup) *AuthRoute {
	return &AuthRoute{
		rg: rg,
	}
}

func (ar *AuthRoute) RegisterRoutes(es *elasticsearch.Client) {
}

func (ar *AuthRoute) search()
