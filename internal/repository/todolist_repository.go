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

func (r *TodolistRepository) Create(ctx context.Context, arg *model.TodoList) error {
	trx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Ensure rollback if anything fails
	defer func() {
		if err != nil {
			trx.Rollback()
		}
	}()

	query := `
		INSERT INTO todolists (title, description)
		VALUES ($1, $2)
		RETURNING id, title, description
	`

	err = trx.QueryRowContext(
		ctx, 
		query, 
		arg.Title, 
		arg.Description,
	).Scan(
		&arg.ID,
		&arg.Title,
		&arg.Description,
	)

	if err != nil {
		return err
	}

	return trx.Commit()
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

func (r *TodolistRepository) Update(ctx context.Context, arg *model.TodoList, id int64) (*model.TodoList, error) {
	trx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			trx.Rollback()
		}
	}()

	query := `
		UPDATE todolists
		SET 
			title = COALESCE($1, title), 
			description = COALESCE($2, description)
			updated_at = NOW()
		WHERE id = $3
	`

	var t model.TodoList

	err = trx.QueryRowContext(
		ctx, 
		query, 
		arg.Title, 
		arg.Description,
		id,
	).Scan(
		&t.Title, 
		&t.Description,
		&t.ID,
	)

	if err != nil {
		return nil, err
	}

	trx.Commit()
	return &t, nil
}

func (r *TodolistRepository) Delete(ctx context.Context, id int64) error {
	trx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			trx.Rollback()
		}
	}()

	query := `
		DELETE FROM todolists
		WHERE id = $1
	`

	_, err = trx.ExecContext(
		ctx, 
		query, 
		id,
	)

	if err != nil {
		return nil
	}

	trx.Commit()

	return nil
}
