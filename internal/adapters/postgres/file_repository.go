package postgres

import (
	"database/sql"
	"wanny-web-services/internal/core/domain"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file *domain.File) error {
	_, err := r.db.Exec("INSERT INTO files (user_id, filename, size) VALUES ($1, $2, $3)", file.UserID, file.Filename, file.Size)
	return err
}

func (r *FileRepository) GetByUserAndFilename(userID int64, filename string) (*domain.File, error) {
	file := &domain.File{}
	err := r.db.QueryRow("SELECT id, user_id, filename, size FROM files WHERE user_id = $1 AND filename = $2", userID, filename).Scan(&file.ID, &file.UserID, &file.Filename, &file.Size)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (r *FileRepository) Delete(fileID int64) error {
	_, err := r.db.Exec("DELETE FROM files WHERE id = $1", fileID)
	return err
}
