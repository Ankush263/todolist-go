package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Ankush263/todolist/internal/model"
)


type UserRepository struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, arg *model.User) error {
	query := `
		INSERT INTO	users (firstname, lastname, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id, firstname, lastname, email, password
	`
	
	return r.db.QueryRowContext(
		ctx, 
		query, 
		&arg.FirstName, 
		&arg.LastName, 
		&arg.Email, 
		&arg.Password,
	).Scan(
		&arg.ID,
		&arg.FirstName,
		&arg.LastName,
		&arg.Email,
		&arg.Password,
	)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, firstname, lastname, email
		FROM users
		WHERE email = $1
	`

	var u model.User

	err := r.db.QueryRowContext(
		ctx, 
		query, 
		email,
	).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("Invalid Creadentials")
	}
	return &u, err
}
