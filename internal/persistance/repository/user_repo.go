package repository

import (
	"context"

	"github.com/echo-webkom/goat/internal/db"
	"github.com/echo-webkom/goat/internal/domain"
)

type UserRepo interface {
	Get(id string) (*domain.User, error)
}

type happeningRepo struct {
	db *db.DB
}

func NewHappeningRepo(db *db.DB) *happeningRepo {
	return &happeningRepo{db: db}
}

func (r *happeningRepo) Get(id string) (*domain.User, error) {
	sql := `SELECT * FROM users WHERE id = $1`
	row := r.db.QueryRow(context.TODO(), sql, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.EmailVerified, &user.Image, &user.AlternativeEmail, &user.DegreeID, &user.Year, &user.Type, &user.IsBanned, &user.BannedFromStrike)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
