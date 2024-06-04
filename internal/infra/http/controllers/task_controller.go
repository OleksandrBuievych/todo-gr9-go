package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		task.UserId = user.Id
		task.Status = domain.New
		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Created(w, tDto)
	}
}

func (c TaskController) GetForUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		tasks, err := c.taskService.GetForUser(user.Id)
		if err != nil {
			log.Printf("TaskController -> GetForUser: %s", err)
			InternalServerError(w, err)
			return
		}

		var tasksDto resources.TasksDto
		tasksDto = tasksDto.DomainToDtoCollection(tasks)
		Success(w, tasksDto)
	}
}

func (c TaskController) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Printf("TaskController -> GetByID: %s", err)
			BadRequest(w, err)
			return
		}

		task, err := c.taskService.GetByID(id)
		if err != nil {
			log.Printf("TaskController -> GetByID: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Success(w, tDto)
	}
}

func (c TaskController) DeleteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Printf("TaskController -> DeleteByID: %s", err)
			BadRequest(w, err)
			return
		}

		err = c.taskService.DeleteByID(id)
		if err != nil {
			log.Printf("TaskController -> DeleteByID: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, "Завдання успішно видалено!")
	}
}

func (c TaskController) UpdateStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Printf("TaskController -> UpdateStatus: %s", err)
			BadRequest(w, err)
			return
		}

		var statusUpdateReq struct {
			Status domain.TaskStatus `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&statusUpdateReq); err != nil {
			log.Printf("TaskController -> UpdateStatus: %s", err)
			BadRequest(w, err)
			return
		}

		if !isValidTaskStatus(statusUpdateReq.Status) {
			err := fmt.Errorf("Не дійсний статус: %s. Використайте будь-ласка: NEW or IN_PROGRESS or DONE", statusUpdateReq.Status)
			log.Printf("TaskController -> UpdateStatus: %s", err)
			BadRequest(w, err)
			return
		}

		err = c.taskService.UpdateStatus(id, statusUpdateReq.Status)
		if err != nil {
			log.Printf("TaskController -> UpdateStatus: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, "Статус успішно змінено!")
	}
}

func isValidTaskStatus(status domain.TaskStatus) bool {
	switch status {
	case domain.New, domain.InProgress, domain.Done:
		return true
	default:
		return false
	}
}
