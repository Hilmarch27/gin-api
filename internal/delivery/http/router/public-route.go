package router

import (
	"github.com/Hilmarch27/gin-api/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

type PublicRouter  struct {
	authHandler *handler.AuthHandler
	jwtSecret   string
}

func NewPublicRouter(authHandler *handler.AuthHandler, jwtSecret string) *PublicRouter {
	return &PublicRouter{
		authHandler: authHandler,
		jwtSecret:   jwtSecret,
	}
}

func (r *PublicRouter) Setup(engine *gin.Engine) {
	// Public routes
	auth := engine.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/refresh", r.authHandler.RefreshToken)
	}
}