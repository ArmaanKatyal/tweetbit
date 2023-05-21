package services

import "github.com/gin-gonic/gin"

// GinService is a struct that contains a pointer to a gin.Engine
type GinService struct {
	Engine *gin.Engine
}

// NewGinService is a function that returns a pointer to a GinService
func NewGinService() *GinService {
	return &GinService{
		Engine: gin.Default(),
	}
}

// Run is a method that runs the gin.Engine
func (g *GinService) Run() {
	err := g.Engine.Run(":3000")
	if err != nil {
		return 
	}
}
