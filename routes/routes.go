package routes

import (
	"TaskManagerApi/handler"
	"TaskManagerApi/middleware"
	"net/http"
)

func RegisterRoutes(
	mux *http.ServeMux,
	taskHandler *handler.TaskHandler,
	userHandler *handler.UserHandler,
) {
	mux.HandleFunc("/api/auth/register", userHandler.Register)
	mux.HandleFunc("/api/auth/login", userHandler.Login)

	mux.Handle("/api/tasks", middleware.JWT(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				taskHandler.CreateTask(w, r)
			case http.MethodGet:
				taskHandler.GetAllTasks(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	))

	mux.Handle("/api/tasks/", middleware.JWT(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				taskHandler.GetTaskById(w, r)
			case http.MethodPut:
				taskHandler.UpdateTask(w, r)
			case http.MethodDelete:
				taskHandler.DeleteTask(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	))
}
