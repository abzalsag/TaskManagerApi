package dto

type TaskRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
}
