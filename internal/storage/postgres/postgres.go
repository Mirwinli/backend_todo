package postgres

import (
	"bd_service/internal/models"
	"bd_service/internal/storage"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{
		db: db,
	}
}

func (s Storage) CreateTask(ctx context.Context, title string, description string, uid int64) error {
	const op = "storage.CreateTask"
	sql := "INSERT INTO tasks (title,description,uid) VALUES ($1,$2,$3)"

	teg, err := s.db.Exec(ctx, sql, title, description, uid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if teg.RowsAffected() == 0 {
		return fmt.Errorf("%s:%w", op, storage.ErrInternal)
	}
	return nil
}

func (s Storage) DeleteTask(ctx context.Context, title string, uid int64) error {
	const op = "storage.DeleteTask"

	sql := "DELETE FROM tasks WHERE title=$1 AND uid=$2"
	teg, err := s.db.Exec(ctx, sql, title, uid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if teg.RowsAffected() == 0 {
		return fmt.Errorf("%s:%w", op, storage.ErrNotFound)
	}
	return nil
}

func (s Storage) DoneTask(ctx context.Context, title string, uid int64) error {
	sql := "UPDATE tasks SET is_done=TRUE WHERE title=$1 AND uid=$2"

	teg, err := s.db.Exec(ctx, sql, title, uid)
	if err != nil {
		return fmt.Errorf("%s: %w", sql, err)
	}
	if teg.RowsAffected() == 0 {
		return fmt.Errorf("%s:%w", sql, storage.ErrNotFound)
	}
	return nil
}

func (s Storage) ListTasks(ctx context.Context, uid int64) ([]models.Task, error) {
	sql := "SELECT title,description,task_id,created_at,done_at,duration,is_done FROM tasks WHERE uid=$1"

	rows, err := s.db.Query(ctx, sql, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s:%w", storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", sql, err)
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task

		if err = rows.Scan(
			&task.Title,
			&task.Description,
			&task.TaskId,
			&task.CreatedAt,
			&task.DoneAt,
			&task.Duration,
			&task.IsDone,
		); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
