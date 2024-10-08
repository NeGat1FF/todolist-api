package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NeGat1FF/todolist-api/internal/models"
	"github.com/NeGat1FF/todolist-api/internal/service"
)

type TaskHandler struct {
	ser *service.Service
}

func NewTaskHandler(ser *service.Service) *TaskHandler {
	return &TaskHandler{ser}
}

func InternalError(rw http.ResponseWriter) {
	http.Error(rw, "internal server error", http.StatusInternalServerError)
}

type AddTaskRequest struct {
	Title       string
	Description string
}

type UpdateTaskRequest struct {
	Title       *string
	Description *string
}

type GetTasksResponse struct {
	Data  []models.Task
	Page  int
	Limit int
	Total int
}

// AddTask godoc
//
//	@Summary		Add a new task
//	@Description	Add a new task
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			task	body		AddTaskRequest	true	"Task object that needs to be added"
//	@Success		202		{object}	models.Task
//	@Router			/tasks [post]
func (th *TaskHandler) AddTask(rw http.ResponseWriter, r *http.Request) {
	task := r.Context().Value(models.TaskKey{}).(models.Task)
	user_id := r.Context().Value(models.UserIDKey{}).(int)

	task.UserID = user_id

	task, err := th.ser.AddTask(r.Context(), task)
	if err != nil {
		InternalError(rw)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusAccepted)
	json.NewEncoder(rw).Encode(task)
}

// GetTasks godoc
//
//	@Summary		Get all tasks
//	@Description	Get all tasks
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int						false	"Page number"
//	@Param			limit	query		int						false	"Limit number"
//	@Success		200		{object}	GetTasksResponse
//	@Router			/tasks [get]
func (th *TaskHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10

	if p := r.URL.Query().Get("page"); p != "" {
		if pVal, err := strconv.Atoi(p); err == nil && pVal > 0 {
			page = pVal
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if lVal, err := strconv.Atoi(l); err == nil && lVal > 0 {
			limit = lVal
		}
	}

	userID := r.Context().Value(models.UserIDKey{}).(int)

	// Fetch tasks
	tasks, err := th.ser.GetTasks(r.Context(), userID, page, limit)
	if err != nil {
		InternalError(rw)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	response := map[string]any{
		"data":  tasks,
		"page":  page,
		"limit": limit,
		"total": len(tasks),
	}

	json.NewEncoder(rw).Encode(response)
}

// UpdateTask godoc
//
//	@Summary		Update a task
//	@Description	Update a task
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Task ID"
//	@Param			task	body		UpdateTaskRequest	true	"Task object that needs to be updated"
//	@Success		202		{object}	models.Task
//	@Router			/tasks/{id} [put]
func (th *TaskHandler) UpdateTask(rw http.ResponseWriter, r *http.Request) {
	task_id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(rw, "incorrect id", http.StatusBadGateway)
		return
	}
	task := r.Context().Value(models.TaskKey{}).(models.Task)
	user_id := r.Context().Value(models.UserIDKey{}).(int)

	task, err = th.ser.UpdateTask(r.Context(), task, task_id, user_id)
	if err != nil {
		InternalError(rw)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusAccepted)
	json.NewEncoder(rw).Encode(task)
}

// DeleteTask godoc
//
//	@Summary		Delete a task
//	@Description	Delete a task
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//	@Success		204	{object}	string
//	@Router			/tasks/{id} [delete]
func (th *TaskHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	task_id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(rw, "incorrect id", http.StatusBadRequest)
		return
	}
	user_id := r.Context().Value(models.UserIDKey{}).(int)

	err = th.ser.DeleteTask(r.Context(), task_id, user_id)
	if err != nil {
		InternalError(rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
