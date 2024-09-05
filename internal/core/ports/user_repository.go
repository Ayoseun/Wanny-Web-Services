package ports

import "wanny-web-services/internal/core/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
	UpdateUsage(userID int64, usage int64) error
}
