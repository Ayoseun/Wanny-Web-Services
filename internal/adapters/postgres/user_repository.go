package postgres

import (
	"database/sql"
	"wanny-web-services/internal/core/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	return err
}

func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	err := r.db.QueryRow("SELECT id, username, password, usage FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password, &user.Usage)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUsage(userID int64, usage int64) error {
	_, err := r.db.Exec("UPDATE users SET usage = usage + $1 WHERE id = $2", usage, userID)
	return err
}
