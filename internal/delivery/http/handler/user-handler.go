package handler

import (
	"net/http"

	"github.com/Hilmarch27/gin-api/internal/domain"
	"github.com/Hilmarch27/gin-api/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
    // Parsing input dari request body
    var req domain.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Panggil usecase untuk memproses registrasi
    if err := h.authUsecase.Register(&req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Response sukses
    c.JSON(http.StatusCreated, gin.H{
        "status":  "success",
        "message": "user registered successfully",
    })
}


func (h *AuthHandler) Login(c *gin.Context) {
    var req domain.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
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

func (h *AuthHandler) GetUserByID(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists || user == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "UserId not found"})
        return
    }

    // Type assertion ke *domain.User
    userObj, ok := user.(*domain.User)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user object"})
        return
    }

    userId := userObj.ID

    userResponse, err := h.authUsecase.GetUserByID(userId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data":   userResponse,
    })
}

func (h *AuthHandler) Update(c *gin.Context) {
    // Ambil ID dari parameter URL
    userId, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Bind data JSON ke UpdateRequest tanpa ID
    var req domain.UpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "detail": err.Error()})
        return
    }

    // Tambahkan ID dari URL ke objek request
    req.ID = userId

    // Panggil usecase untuk update user
    if err := h.authUsecase.UpdateUser(&req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Kirimkan response sukses
    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "User updated successfully",
    })
}

func (h *AuthHandler) Delete(c *gin.Context) {
    // Ambil ID dari parameter URL
    userId, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Panggil usecase untuk delete user
    if err := h.authUsecase.DeleteUser(userId); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Kirimkan response sukses
    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "User deleted successfully",
    })
}