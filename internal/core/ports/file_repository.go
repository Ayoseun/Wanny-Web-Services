package ports

import "wanny-web-services/internal/core/domain"

type FileRepository interface {
	Create(file *domain.File) error
	GetByUserAndFilename(userID int64, filename string) (*domain.File, error)
	Delete(fileID int64) error
}
