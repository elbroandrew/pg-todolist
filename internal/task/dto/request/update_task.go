package request

type UpdateTask struct {
    Title       *string `json:"title" validate:"omitempty"`
    Completed   *bool   `json:"completed"`
}