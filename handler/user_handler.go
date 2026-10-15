package handler

import (
	"TaskManagerApi/dto"
	"TaskManagerApi/service"
	"encoding/json"
	"errors"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		service: userService,
	}
}

func (handler *UserHandler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := handler.service.Register(
		r.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		http.Error(
			w,
			"registration failed",
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "user registered successfully",
	})
}

func (handler *UserHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid request",
			http.StatusBadRequest,
		)
		return
	}

	token, err := handler.service.Login(
		r.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(
				w,
				"invalid credentials",
				http.StatusUnauthorized,
			)
			return
		}

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
