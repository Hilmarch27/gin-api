package repository

import (
	"github.com/Hilmarch27/gin-api/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindById(id uuid.UUID) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uuid.UUID) error
}