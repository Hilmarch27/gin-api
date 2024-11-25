package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Hilmarch27/gin-api/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func AuthenticationMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil access_token dari cookie
		cookie, err := c.Cookie("access_token")
		fmt.Println("Access Token:", cookie)
		if err != nil {
			// Jika tidak ada cookie access_token, lanjutkan
			c.Next()
			return
		}

		// Verifikasi token
		token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
			// Pastikan token menggunakan signing method yang benar
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			// Jika token tidak valid atau expired, lanjutkan tanpa mengatur user
			c.Next()
			return
		}

		// Extract claims dan simpan di context Gin
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Ambil userId dari claims
			userIDStr, ok := claims["userId"].(string)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return
			}

			// Parse userId menjadi uuid.UUID
			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid userId format"})
				c.Abort()
				return
			}

			// Ambil role dari claims
			role, okRole := claims["role"].(string)
			if !okRole {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return
			}

			// Set user info ke context Gin
			user := &domain.User{
				ID:   userID,
				Role: role,
			}
			c.Set("user", user)
		}

		c.Next()
	}
}

func RequireCredentials() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Periksa apakah user sudah ada di konteks
		user, exists := c.Get("user")
		if !exists || user == nil {
			// Jika tidak ada user, kirim response Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Lanjutkan request jika user terautentikasi
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Periksa apakah ada user di konteks
		user, exists := c.Get("user")
		if !exists || user == nil {
			// Jika user tidak ada, kirim response Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		// Cek apakah user memiliki role 'admin'
		if u, ok := user.(*domain.User); ok && u.Role != "admin" {
			// Jika bukan admin, kirim response Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized - Admin access required"})
			c.Abort()
			return
		}

		// Lanjutkan request jika user memiliki akses admin
		c.Next()
	}
}