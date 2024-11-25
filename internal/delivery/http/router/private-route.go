package router

import (
	"github.com/Hilmarch27/gin-api/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
	jwtSecret string
}

func NewApiRouter(jwtSecret string) *ApiRouter {
	return &ApiRouter{
		jwtSecret: jwtSecret,
	}
}

func (r *ApiRouter) Setup(engine *gin.Engine) {
    api := engine.Group("/api")
    api.Use(middleware.RequireCredentials())
    {
        api.GET("/profile", func(c *gin.Context) {
            currUser, _ := c.Get("user")
            c.JSON(200, gin.H{
                "status": "success",
                "data": gin.H{
                    "user": currUser,
                },
            })
        })
    }
}