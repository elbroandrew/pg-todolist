package request

type CreateTask struct {
    Title       string `json:"title" validate:"required"`
}