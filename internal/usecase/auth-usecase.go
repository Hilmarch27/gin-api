package usecase

import (
	"errors"
	"time"

	"github.com/Hilmarch27/gin-api/internal/domain"
	"github.com/Hilmarch27/gin-api/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
    Register(req *domain.RegisterRequest) error
    Login(req *domain.LoginRequest) (string, string, error)
    RefreshToken(refreshToken string) (string, string, error)
}

type authUsecase struct {
    userRepo    repository.UserRepository
    jwtSecret   []byte
    tokenExpiry time.Duration
}

func NewAuthUsecase(ur repository.UserRepository, secret string, expiry time.Duration) AuthUsecase {
    return &authUsecase{
        userRepo:    ur,
        jwtSecret:   []byte(secret),
        tokenExpiry: expiry,
    }
}



func (u *authUsecase) Register(req *domain.RegisterRequest) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := &domain.User{
        Name:     req.Name,
        Email:    req.Email,
        Role:     req.Role,
        Password: string(hashedPassword),
    }


    return u.userRepo.Create(user)
}

func (u *authUsecase) generateTokens(userID uuid.UUID, userRole string) (string, string, error) {
    // Generate Access Token
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userId": userID,
        "role":   userRole,
        "exp":    time.Now().Add(u.tokenExpiry).Unix(),
    })
    accessTokenString, err := accessToken.SignedString(u.jwtSecret)
    if err != nil {
        return "", "", err
    }

    // Generate Refresh Token
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userId": userID,
        "exp":    time.Now().Add(7 * 24 * time.Hour).Unix(), // Refresh token valid for 1 week
    })
    refreshTokenString, err := refreshToken.SignedString(u.jwtSecret)
    if err != nil {
        return "", "", err
    }

    return accessTokenString, refreshTokenString, nil
}

func (u *authUsecase) Login(req *domain.LoginRequest) (string, string, error) {
    user, err := u.userRepo.FindByEmail(req.Email)
    if err != nil {
        return "", "", errors.New("invalid credentials")
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return "", "", errors.New("invalid credentials")
    }

    return u.generateTokens(user.ID, user.Role)
}

func (u *authUsecase) RefreshToken(refreshToken string) (string, string, error) {
    // Parse the refresh token
    token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return u.jwtSecret, nil
    })
    if err != nil || !token.Valid {
        return "", "", errors.New("invalid refresh token")
    }

    // Extract claims from the refresh token
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return "", "", errors.New("invalid refresh token claims")
    }

    userIDStr, ok := claims["userId"].(string)
    if !ok {
        return "", "", errors.New("userId not found in refresh token")
    }

    // Parse the string userID into uuid.UUID
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return "", "", errors.New("invalid userId format")
    }

    // Find the user by ID
    user, err := u.userRepo.FindById(userID)
    if err != nil {
        return "", "", errors.New("user not found")
    }

    // Generate new access token and refresh token
    accessToken, newRefreshToken, err := u.generateTokens(user.ID, user.Role)
    if err != nil {
        return "", "", err
    }

    return accessToken, newRefreshToken, nil
}