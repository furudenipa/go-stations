package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	sqlResult, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}

	id, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID:          id,
		Subject:     subject,
		Description: description,
	}

	row := s.db.QueryRowContext(ctx, confirm, id)
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	var rows *sql.Rows
	if prevID == 0 {
		stmt, err := s.db.PrepareContext(ctx, read)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		rows, err = stmt.QueryContext(ctx, size)
		if err != nil {
			return nil, err
		}
	} else {
		stmt, err := s.db.PrepareContext(ctx, readWithID)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()

		rows, err = stmt.QueryContext(ctx, prevID, size)
		if err != nil {
			return nil, err
		}
	}

	todos := make([]*model.TODO, 0, size)
	for rows.Next() {
		todo := &model.TODO{}
		if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	sqlResult, err := stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}

	numAffectedRow, err := sqlResult.RowsAffected()
	if err != nil {
		return nil, err
	}
	if numAffectedRow == 0 {
		return nil, &model.ErrNotFound{}
	}

	todo := &model.TODO{
		ID:          id,
		Subject:     subject,
		Description: description,
	}

	row := s.db.QueryRowContext(ctx, confirm, id)
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`
	var delete string
	if numIDs := len(ids); numIDs == 0 {
		return nil
	} else {
		delete = fmt.Sprintf(deleteFmt, strings.Repeat(",?", numIDs-1))
	}

	stmt, err := s.db.PrepareContext(ctx, delete)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var args []interface{}
	for _, id := range ids {
		args = append(args, id)
	}
	sqlResult, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}
	if numAffectedRows, err := sqlResult.RowsAffected(); err != nil {
		return err
	} else if numAffectedRows == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}
