package service

import (
	"pg-todolist/internal/task/repository"
	"pg-todolist/pkg/app_errors"
)

type SyncService struct {
	taskRepo repository.TaskRepository
	cacheRepository repository.CacheRepository
}

func NewSyncService(taskRepo repository.TaskRepository, cacheRepository repository.CacheRepository) *SyncService {
	return &SyncService{
		taskRepo: taskRepo,
		cacheRepository: cacheRepository,
	}
}

// SyncUserTasks синхронизирует задачи пользователя между БД и кешем
func (s *SyncService) SyncUserTasks(userID uint) *app_errors.AppError {
	// Получаем актуальные задачи из БД
	tasks, err := s.taskRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	// Обновляем кеш
	if cacheErr := s.cacheRepository.SetUserTasks(userID, tasks); cacheErr != nil {
		return app_errors.ErrInternalServer
	}

	return nil
}

// SyncCompletedTasks помечает завершенные задачи в кеше
func (s *SyncService) SyncCompletedTasks(taskIDs []uint) *app_errors.AppError {
	for _, taskID := range taskIDs {
		if err := s.cacheRepository.MarkCompleted(taskID); err != nil {
			return app_errors.ErrInternalServer
		}
	}
	return nil
}

// SyncDeletedTasks помечает удаленные задачи в кеше
func (s *SyncService) SyncDeletedTasks(taskIDs []uint) *app_errors.AppError {
	for _, taskID := range taskIDs {
		if err := s.cacheRepository.MarkDeleted(taskID); err != nil {
			return app_errors.ErrInternalServer
		}
	}
	return nil
}

// FullSync выполняет полную синхронизацию для пользователя
func (s *SyncService) FullSync(userID uint, completedTaskIDs, deletedTaskIDs []uint) *app_errors.AppError {
	// Синхронизируем все задачи пользователя
	if err := s.SyncUserTasks(userID); err != nil {
		return err
	}

	// Синхронизируем завершенные задачи
	if err := s.SyncCompletedTasks(completedTaskIDs); err != nil {
		return err
	}

	// Синхронизируем удаленные задачи
	if err := s.SyncDeletedTasks(deletedTaskIDs); err != nil {
		return err
	}

	return nil
}