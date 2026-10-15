package repository

import (
	"TaskManagerApi/model"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) CreateUser(user *model.User) error {
	_, err := repo.db.Exec(
		`INSERT INTO users (email, password_hash)
		 VALUES ($1, $2)`,
		user.Email,
		user.Password,
	)

	return err
}

func (repo *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User

	err := repo.db.QueryRow(
		`SELECT id, email, password_hash
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}

		return nil, err
	}

	return &user, nil
}
