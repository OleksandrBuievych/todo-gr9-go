package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type TaskService interface {
	Save(t domain.Task) (domain.Task, error)
	GetForUser(uId uint64) ([]domain.Task, error)
	GetByID(id uint64) (domain.Task, error)
	DeleteByID(id uint64) error
	UpdateStatus(id uint64, status domain.TaskStatus) error
}

type taskService struct {
	taskRepo database.TaskRepository
}

func NewTaskService(tr database.TaskRepository) TaskService {
	return taskService{
		taskRepo: tr,
	}
}

func (s taskService) Save(t domain.Task) (domain.Task, error) {
	task, err := s.taskRepo.Save(t)
	if err != nil {
		log.Printf("TaskService -> Save: %s", err)
		return domain.Task{}, err
	}
	return task, nil
}

func (s taskService) GetForUser(uId uint64) ([]domain.Task, error) {
	tasks, err := s.taskRepo.GetByUserId(uId)
	if err != nil {
		log.Printf("TaskService -> GetForUser: %s", err)
		return nil, err
	}
	return tasks, nil
}

func (s taskService) GetByID(id uint64) (domain.Task, error) {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		log.Printf("TaskService -> GetByID: %s", err)
		return domain.Task{}, err
	}
	return task, nil
}

func (s taskService) DeleteByID(id uint64) error {
	err := s.taskRepo.DeleteByID(id)
	if err != nil {
		log.Printf("TaskService -> DeleteByID: %s", err)
		return err
	}
	return nil
}

func (s taskService) UpdateStatus(id uint64, status domain.TaskStatus) error {
	err := s.taskRepo.UpdateStatus(id, status)
	if err != nil {
		log.Printf("TaskService -> UpdateStatus: %s", err)
		return err
	}
	return nil
}
