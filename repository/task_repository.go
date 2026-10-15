package repository

import (
	"TaskManagerApi/model"
	"context"
	"database/sql"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}
func (repo *TaskRepository) CreateTask(
	ctx context.Context,
	task *model.Task,
) error {
	query := `
		INSERT INTO tasks
			(user_id, title, description, completed)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	return repo.db.QueryRowContext(
		ctx,
		query,
		task.UserID,
		task.Title,
		task.Description,
		task.Completed,
	).Scan(&task.ID)
}
func (repo *TaskRepository) GetAllTasks(
	ctx context.Context,
	userID int,
) ([]model.Task, error) {

	rows, err := repo.db.QueryContext(
		ctx,
		`SELECT id, user_id, title, description, completed
		 FROM tasks
		 WHERE user_id = $1
		 ORDER BY id`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task

	for rows.Next() {
		var task model.Task

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.Completed,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (repo *TaskRepository) GetTaskById(
	ctx context.Context,
	taskID int,
	userID int,
) (model.Task, error) {

	var task model.Task

	err := repo.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, title, description, completed
		 FROM tasks
		 WHERE id = $1 AND user_id = $2`,
		taskID,
		userID,
	).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Completed,
	)

	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (repo *TaskRepository) UpdateTask(
	ctx context.Context,
	task model.Task,
) error {

	result, err := repo.db.ExecContext(
		ctx,
		`UPDATE tasks
		 SET title = $1,
		     description = $2,
		     completed = $3
		 WHERE id = $4 AND user_id = $5`,
		task.Title,
		task.Description,
		task.Completed,
		task.ID,
		task.UserID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (repo *TaskRepository) DeleteTask(
	ctx context.Context,
	taskID int,
	userID int,
) error {

	result, err := repo.db.ExecContext(
		ctx,
		`DELETE FROM tasks
		 WHERE id = $1 AND user_id = $2`,
		taskID,
		userID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
