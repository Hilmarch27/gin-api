package router

import (
	"github.com/Hilmarch27/gin-api/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	authHandler *handler.AuthHandler
	jwtSecret   string
}

func NewAuthRouter(authHandler *handler.AuthHandler, jwtSecret string) *AuthRouter {
	return &AuthRouter{
		authHandler: authHandler,
		jwtSecret:   jwtSecret,
	}
}

func (r *AuthRouter) Setup(engine *gin.Engine) {
	// Public routes
	auth := engine.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/refresh", r.authHandler.RefreshToken)
	}
}