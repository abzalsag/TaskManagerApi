package main

import (
	"encoding/json"
	"testing"
)

func TestTaskCreatedEventUnmarshal(t *testing.T) {
	data := []byte(`{
		"task_id": 10,
		"user_id": 1,
		"title": "Test task",
		"description": "Testing event"
	}`)

	var event TaskCreatedEvent

	err := json.Unmarshal(data, &event)
	if err != nil {
		t.Fatalf("failed to unmarshal event: %v", err)
	}

	if event.TaskID != 10 {
		t.Errorf("expected task ID 10, got %d", event.TaskID)
	}

	if event.UserID != 1 {
		t.Errorf("expected user ID 1, got %d", event.UserID)
	}

	if event.Title != "Test task" {
		t.Errorf("expected title 'Test task', got %s", event.Title)
	}

	if event.Description != "Testing event" {
		t.Errorf(
			"expected description 'Testing event', got %s",
			event.Description,
		)
	}
}

func TestHandleTaskCreatedEventInvalidJSON(t *testing.T) {
	data := []byte(`invalid json`)

	err := handleTaskCreatedEvent(data)

	if err == nil {
		t.Error("expected error, got nil")
	}
}
