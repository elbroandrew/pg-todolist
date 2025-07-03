package response

type TaskList struct {
    Tasks  []Task `json:"tasks"`
    Total  int    `json:"total"`
    Limit  int    `json:"limit"`
    Offset int    `json:"offset"`
}