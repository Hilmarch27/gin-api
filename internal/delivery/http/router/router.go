package router

import (
    "github.com/gin-gonic/gin"
    "github.com/Hilmarch27/gin-api/internal/delivery/http/middleware"
)

type Router struct {
    engine     *gin.Engine
    auth       *PublicRouter
    api        *ApiRouter
    jwtSecret  []byte
}

func NewRouter(engine *gin.Engine, authRouter *PublicRouter, apiRouter *ApiRouter, jwtSecret []byte) *Router {
    return &Router{
        engine: engine,
        auth:   authRouter,
        api:    apiRouter,
        jwtSecret: jwtSecret,
    }
}

func (r *Router) SetupRoutes() {
    // Setup global middlewares
    r.engine.Use(gin.Logger())
    r.engine.Use(gin.Recovery())
    
    // Add authentication middleware globally
    r.engine.Use(middleware.AuthenticationMiddleware(r.jwtSecret))

    // Setup route groups
    // Auth routes (public)
    r.auth.Setup(r.engine)
    
    // API routes (private)
    r.api.Setup(r.engine)
}
