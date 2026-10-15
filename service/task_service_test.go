package service

import (
	"TaskManagerApi/model"
	"context"
	"errors"
	"fmt"
	"testing"

	"TaskManagerApi/kafka"
)

type mockTaskRepository struct {
	task *model.Task
	err  error
}

func (m *mockTaskRepository) CreateTask(
	ctx context.Context,
	task *model.Task,
) error {
	if m.err != nil {
		return m.err
	}

	task.ID = 100
	m.task = task

	return nil
}

func (m *mockTaskRepository) GetAllTasks(
	ctx context.Context,
	userID int,
) ([]model.Task, error) {
	return nil, nil
}

func (m *mockTaskRepository) GetTaskById(
	ctx context.Context,
	taskID int,
	userID int,
) (model.Task, error) {
	return model.Task{}, nil
}

func (m *mockTaskRepository) UpdateTask(
	ctx context.Context,
	task model.Task,
) error {
	return nil
}

func (m *mockTaskRepository) DeleteTask(
	ctx context.Context,
	taskID int,
	userID int,
) error {
	return nil
}

type mockTaskEventPublisher struct {
	called bool
	event  kafka.TaskCreatedEvent
	err    error
}

func (m *mockTaskEventPublisher) PublishTaskCreated(
	ctx context.Context,
	event kafka.TaskCreatedEvent,
) error {
	if m.err != nil {
		return m.err
	}

	m.called = true
	m.event = event

	return nil
}

func TestTaskServiceCreateTaskPublishesEvent(t *testing.T) {

	publisher := &mockTaskEventPublisher{}

	event := kafka.TaskCreatedEvent{
		TaskID:      1,
		UserID:      10,
		Title:       "Test task",
		Description: "Test description",
	}

	err := publisher.PublishTaskCreated(
		context.Background(),
		event,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !publisher.called {
		t.Fatal("expected publisher to be called")
	}

	if publisher.event.TaskID != 1 {
		t.Errorf(
			"expected task ID 1, got %d",
			publisher.event.TaskID,
		)
	}

	if publisher.event.UserID != 10 {
		t.Errorf(
			"expected user ID 10, got %d",
			publisher.event.UserID,
		)
	}

	if publisher.event.Title != "Test task" {
		t.Errorf(
			"expected title 'Test task', got %s",
			publisher.event.Title,
		)
	}
}

func TestTaskServiceCreateTask(t *testing.T) {
	repo := &mockTaskRepository{}
	publisher := &mockTaskEventPublisher{}

	service := NewTaskService(repo, publisher)

	task := &model.Task{
		UserID:      10,
		Title:       "Test task",
		Description: "Testing service",
		Completed:   false,
	}

	err := service.CreateTask(
		context.Background(),
		task,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != 100 {
		t.Errorf(
			"expected task ID 100, got %d",
			task.ID,
		)
	}

	if !publisher.called {
		t.Fatal("expected Kafka event to be published")
	}

	if publisher.event.TaskID != 100 {
		t.Errorf(
			"expected event task ID 100, got %d",
			publisher.event.TaskID,
		)
	}

	if publisher.event.UserID != 10 {
		t.Errorf(
			"expected event user ID 10, got %d",
			publisher.event.UserID,
		)
	}

	if publisher.event.Title != "Test task" {
		t.Errorf(
			"expected event title 'Test task', got %s",
			publisher.event.Title,
		)
	}
}

func TestTaskServiceCreateTaskInvalidTask(t *testing.T) {
	repo := &mockTaskRepository{}
	publisher := &mockTaskEventPublisher{}

	service := NewTaskService(repo, publisher)

	task := &model.Task{
		UserID:      10,
		Title:       "",
		Description: "Invalid task",
	}

	err := service.CreateTask(
		context.Background(),
		task,
	)

	if !errors.Is(err, ErrInvalidTask) {
		t.Fatalf(
			"expected ErrInvalidTask, got %v",
			err,
		)
	}

	if repo.task != nil {
		t.Error("repository should not be called")
	}

	if publisher.called {
		t.Error("publisher should not be called")
	}
}

func TestTaskServiceCreateTaskRepositoryError(t *testing.T) {
	repositoryError := fmt.Errorf("database error")

	repo := &mockTaskRepository{
		err: repositoryError,
	}

	publisher := &mockTaskEventPublisher{}

	service := NewTaskService(repo, publisher)

	task := &model.Task{
		UserID:      10,
		Title:       "Test task",
		Description: "Testing repository error",
	}

	err := service.CreateTask(
		context.Background(),
		task,
	)

	if !errors.Is(err, repositoryError) {
		t.Fatalf(
			"expected repository error, got %v",
			err,
		)
	}

	if publisher.called {
		t.Error("publisher should not be called")
	}
}

func TestTaskServiceCreateTaskKafkaError(t *testing.T) {
	kafkaError := fmt.Errorf("kafka error")

	repo := &mockTaskRepository{}

	publisher := &mockTaskEventPublisher{
		err: kafkaError,
	}

	service := NewTaskService(repo, publisher)

	task := &model.Task{
		UserID:      10,
		Title:       "Test task",
		Description: "Testing Kafka error",
	}

	err := service.CreateTask(
		context.Background(),
		task,
	)

	if !errors.Is(err, kafkaError) {
		t.Fatalf(
			"expected Kafka error, got %v",
			err,
		)
	}

	if task.ID != 100 {
		t.Errorf(
			"expected task ID 100, got %d",
			task.ID,
		)
	}
}
