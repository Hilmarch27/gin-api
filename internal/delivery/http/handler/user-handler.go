package handler

import (
	"net/http"

	"github.com/Hilmarch27/gin-api/internal/domain"
	"github.com/Hilmarch27/gin-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authUsecase usecase.AuthUsecase
}

func NewAuthHandler(au usecase.AuthUsecase) *AuthHandler {
    return &AuthHandler{
        authUsecase: au,
    }
}

func (h *AuthHandler) Register(c *gin.Context) {
    var req domain.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.authUsecase.Register(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "status":  "success",
        "message": "user registered successfully",
    })
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req domain.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    accessToken, refreshToken, err := h.authUsecase.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Set cookies
    c.SetCookie("access_token", accessToken, 3600, "/", "", false, true)  // 1 hour
    c.SetCookie("refresh_token", refreshToken, 604800, "/", "", false, true) // 1 week

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "login successful",
    })
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
    // Ambil refresh token dari cookie
    refreshToken, err := c.Cookie("refresh_token")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
        return
    }

    // Panggil usecase untuk refresh token
    accessToken, newRefreshToken, err := h.authUsecase.RefreshToken(refreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Set cookies baru untuk access token dan refresh token
    c.SetCookie("access_token", accessToken, 3600, "/", "", false, true)
    c.SetCookie("refresh_token", newRefreshToken, 604800, "/", "", false, true)

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "tokens refreshed successfully",
    })
}