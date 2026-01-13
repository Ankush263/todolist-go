package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Ankush263/todolist/internal/model"
)

type TodolistRepository struct {
	db *sql.DB
}

func NewTodolistRepo(db *sql.DB) *TodolistRepository {
	return &TodolistRepository{db: db}
}

func (r *TodolistRepository) Create(ctx context.Context, arg model.TodoList) error {
	query := `
		INSERT INTO todolists (title, description)
		VALUES ($1, $2)
		RETURNING id, title, description
	`

	return r.db.QueryRowContext(
		ctx, 
		query, 
		&arg.Title, 
		&arg.Description,
	).Scan(
		&arg.ID,
		&arg.Title,
		&arg.Description,
	)
}

func (r *TodolistRepository) GetAllTodolistsOfUser(ctx context.Context, userid int64) (*model.TodoList, error) {
	query := `
		SELECT id, title, description 
		FROM todolists
		WHERE created_by = $1
	`

	var todolist model.TodoList

	err := r.db.QueryRowContext(
		ctx, 
		query, 
		userid,
	).Scan(
		&todolist.ID,
		&todolist.CreatedBy,
		&todolist.Title,
		&todolist.Description,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("Invalid Credentials")
	}

	return &todolist, nil
}

func (r *TodolistRepository) GetSingleTodolist(ctx context.Context, id int64) (*model.TodoList, error) {
	query := `
		SELECT id, title, description
		FROM todolist
		WHERE id = $1
	`
	var t model.TodoList

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Description,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("Invalid Credentials")
	}
	return &t, nil
}
