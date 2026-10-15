package kafka

type TaskCreatedEvent struct {
	TaskID      int    `json:"task_id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
