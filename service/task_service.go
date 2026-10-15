package service

import (
	"TaskManagerApi/kafka"
	"TaskManagerApi/model"
	"context"
	"database/sql"
	"errors"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *model.Task) error
	GetAllTasks(ctx context.Context, userID int) ([]model.Task, error)
	GetTaskById(ctx context.Context, taskID int, userID int) (model.Task, error)
	UpdateTask(ctx context.Context, task model.Task) error
	DeleteTask(ctx context.Context, taskID int, userID int) error
}

type TaskEventPublisher interface {
	PublishTaskCreated(
		ctx context.Context,
		event kafka.TaskCreatedEvent,
	) error
}

type TaskService struct {
	repo     TaskRepository
	producer TaskEventPublisher
}

func NewTaskService(
	repo TaskRepository,
	producer TaskEventPublisher,
) *TaskService {
	return &TaskService{
		repo:     repo,
		producer: producer,
	}
}

func (service *TaskService) CreateTask(
	ctx context.Context,
	task *model.Task,
) error {
	if err := validateTask(task); err != nil {
		return err
	}

	if err := service.repo.CreateTask(ctx, task); err != nil {
		return err
	}

	event := kafka.TaskCreatedEvent{
		TaskID:      task.ID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
	}

	if err := service.producer.PublishTaskCreated(ctx, event); err != nil {
		return err
	}

	return nil
}

func (service *TaskService) GetAllTasks(
	ctx context.Context,
	userID int,
) ([]model.Task, error) {
	return service.repo.GetAllTasks(ctx, userID)
}

func (service *TaskService) GetTaskById(
	ctx context.Context,
	taskID int,
	userID int,
) (model.Task, error) {
	task, err := service.repo.GetTaskById(
		ctx,
		taskID,
		userID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, ErrTaskNotFound
	}

	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func (service *TaskService) UpdateTask(
	ctx context.Context,
	task model.Task,
) error {
	if err := validateTask(&task); err != nil {
		return err
	}

	err := service.repo.UpdateTask(ctx, task)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrTaskNotFound
	}

	return err
}

func (service *TaskService) DeleteTask(
	ctx context.Context,
	taskID int,
	userID int,
) error {
	err := service.repo.DeleteTask(
		ctx,
		taskID,
		userID,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrTaskNotFound
	}

	return err
}
