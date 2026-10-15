package handler

import (
	"TaskManagerApi/model"
	"TaskManagerApi/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func getUserID(r *http.Request) (int, bool) {
	value := r.Context().Value("user_id")

	userID, ok := value.(float64)
	if !ok {
		return 0, false
	}

	return int(userID), true
}
func (handler *TaskHandler) CreateTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "user id not found", http.StatusUnauthorized)
		return
	}

	var task model.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	task.UserID = userID

	if err := handler.service.CreateTask(
		r.Context(),
		&task,
	); err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(task)
}
func (handler *TaskHandler) GetAllTasks(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "user id not found", http.StatusUnauthorized)
		return
	}

	tasks, err := handler.service.GetAllTasks(
		r.Context(),
		userID,
	)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(tasks)
}

func (handler *TaskHandler) GetTaskById(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "user id not found", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(
		r.URL.Path,
		"/api/tasks/",
	)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := handler.service.GetTaskById(
		r.Context(),
		id,
		userID,
	)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(task)
}

func (handler *TaskHandler) UpdateTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "user id not found", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(
		r.URL.Path,
		"/api/tasks/",
	)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var task model.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	task.ID = id
	task.UserID = userID

	if err := handler.service.UpdateTask(
		r.Context(),
		task,
	); err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(task)
}

func (handler *TaskHandler) DeleteTask(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := getUserID(r)
	if !ok {
		http.Error(w, "user id not found", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(
		r.URL.Path,
		"/api/tasks/",
	)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := handler.service.DeleteTask(
		r.Context(),
		id,
		userID,
	); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidTask):
		http.Error(
			w,
			"invalid task",
			http.StatusBadRequest,
		)

	case errors.Is(err, service.ErrTaskNotFound):
		http.Error(
			w,
			"task not found",
			http.StatusNotFound,
		)

	default:
		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
	}
}
