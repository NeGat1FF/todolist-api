package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/NeGat1FF/todolist-api/internal/models"
)

func validateTask(r *http.Request, requireBoth bool) (models.Task, error) {
	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		return task, errors.New("failed to parse body")
	}

	if requireBoth {
		if task.Title == "" {
			return task, errors.New("task title is not specified")
		}
		if task.Description == "" {
			return task, errors.New("task description is not specified")
		}
	} else {
		if task.Title == "" && task.Description == "" {
			return task, errors.New("it least one field is required")
		}
	}

	return task, nil
}

func ValidateAddTask(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		task, err := validateTask(r, true)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), models.TaskKey{}, task)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}

func ValidateUpdateTask(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		task, err := validateTask(r, false)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), models.TaskKey{}, task)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	}
}
