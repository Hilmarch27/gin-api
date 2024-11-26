package router

import (
	"github.com/Hilmarch27/gin-api/internal/delivery/http/handler"
	"github.com/Hilmarch27/gin-api/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

type ApiRouter struct {
    authHandler *handler.AuthHandler
	jwtSecret string
}

func NewApiRouter(authHandler *handler.AuthHandler, jwtSecret string) *ApiRouter {
	return &ApiRouter{
        authHandler: authHandler,
		jwtSecret: jwtSecret,
	}
}

func (r *ApiRouter) Setup(engine *gin.Engine) {
    api := engine.Group("/api")
    api.Use(middleware.RequireCredentials())
    {   
        users := api.Group("/users")
        users.GET("", r.authHandler.GetUserByID)
        users.PATCH("/:id", r.authHandler.Update)
        users.DELETE("/:id", r.authHandler.Delete)
    }
    // Tambahkan route admin di sini
    admin := api.Group("/admin")
    admin.Use(middleware.RequireAdmin()) // Tambahkan middleware role admin
    {
        admin.GET("", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "status": "success", 
                "message": "Admin Dashboard",
                "data": gin.H{
                    "total_users": 100,
                    "total_products": 50,
                },
            })
        })
    }
}