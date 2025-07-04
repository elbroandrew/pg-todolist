package sync

import (
	"pg-todolist/internal/task/service"
	"pg-todolist/pkg/app_errors"
)

type SyncWorker struct {
	syncService *service.SyncService
}

func NewSyncWorker(syncService *service.SyncService) *SyncWorker {
	return &SyncWorker{syncService: syncService}
}

// ProcessUserSync обрабатывает синхронизацию для конкретного пользователя
func (w *SyncWorker) ProcessUserSync(userID uint) *app_errors.AppError {
	return w.syncService.SyncUserTasks(userID)
}

// ProcessBatchSync обрабатывает пакетную синхронизацию
func (w *SyncWorker) ProcessBatchSync(userIDs []uint) *app_errors.AppError {
	for _, userID := range userIDs {
		if err := w.ProcessUserSync(userID); err != nil {
			return err
		}
	}
	return nil
}