package dto

// Запрос на принудительную синхронизацию
type SyncRequest struct {
    TaskIDs []string `json:"task_ids"`  // Опционально: IDs задач для синхронизации
}

// Статус синхронизации
type SyncStatus struct {
    SyncedTasks int `json:"synced_tasks"`
    Errors      int `json:"errors"`
}