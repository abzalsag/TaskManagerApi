package service

import (
	"errors"
	"strings"

	"TaskManagerApi/model"
)

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrInvalidTask        = errors.New("invalid task")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func validateTask(task *model.Task) error {
	if strings.TrimSpace(task.Title) == "" {
		return ErrInvalidTask
	}

	return nil
}
