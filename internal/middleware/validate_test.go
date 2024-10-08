package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NeGat1FF/todolist-api/internal/models"
)

func TestValidateAddTask(t *testing.T) {
	next := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}

	testCases := []struct {
		name         string
		task         models.Task
		expectedCode int
	}{
		{
			name: "Title and description",
			task: models.Task{
				Title:       "Test task",
				Description: "This is description of test task",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Only title",
			task: models.Task{
				Title: "Test task",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Only description",
			task: models.Task{
				Description: "This is description of test task",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Not fields",
			task:         models.Task{},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(ValidateAddTask(next))
			defer server.Close()

			data, err := json.Marshal(&test.task)
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest("GET", server.URL, bytes.NewBuffer(data))
			if err != nil {
				t.Error(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != test.expectedCode {
				t.Errorf("expected status code %d, but got: %d", test.expectedCode, resp.StatusCode)
			}
		})
	}
}

func TestValidateUpdateTask(t *testing.T) {
	next := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}

	testCases := []struct {
		name         string
		task         models.Task
		expectedCode int
	}{
		{
			name: "Title and description",
			task: models.Task{
				Title:       "Test task",
				Description: "This is description of test task",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Only title",
			task: models.Task{
				Title: "Test task",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Only description",
			task: models.Task{
				Description: "This is description of test task",
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Not fields",
			task:         models.Task{},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(ValidateUpdateTask(next))
			defer server.Close()

			data, err := json.Marshal(&test.task)
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest("GET", server.URL, bytes.NewBuffer(data))
			if err != nil {
				t.Error(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != test.expectedCode {
				t.Errorf("expected status code %d, but got: %d", test.expectedCode, resp.StatusCode)
			}
		})
	}
}

func TestValidateRegistration(t *testing.T) {
	next := func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}

	testCases := []struct {
		name         string
		user         models.User
		expectedCode int
	}{
		{
			name: "All fields",
			user: models.User{
				Username: "TestUser",
				Email:    "exampleEmail@test.com",
				Password: "Password",
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Incorrect email",
			user: models.User{
				Username: "TestUser",
				Email:    "test.com",
				Password: "Password",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "No Username",
			user: models.User{
				Email:    "exampleEmail@test.com",
				Password: "Password",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "No Email",
			user: models.User{
				Username: "TestUser",
				Password: "Password",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "No Password",
			user: models.User{
				Username: "TestUser",
				Email:    "exampleEmail@test.com",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "No body",
			user:         models.User{},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(ValidateRegistration(next))
			defer server.Close()

			data, err := json.Marshal(&test.user)
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest("GET", server.URL, bytes.NewBuffer(data))
			if err != nil {
				t.Error(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Error(err)
			}

			if resp.StatusCode != test.expectedCode {
				t.Errorf("expected status code %d, but got: %d", test.expectedCode, resp.StatusCode)
			}
		})
	}
}
