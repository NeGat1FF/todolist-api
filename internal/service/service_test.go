package service

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/NeGat1FF/todolist-api/internal/models"
	"github.com/NeGat1FF/todolist-api/mocks"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterUser(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name          string
		inputUser     models.User
		mockSetup     func(userRepoMock *mocks.UserRepositoryInterface)
		expectedError bool
	}{
		{
			name:      "User does not exist",
			inputUser: models.User{Email: "new@example.com"},
			mockSetup: func(userRepoMock *mocks.UserRepositoryInterface) {
				userRepoMock.On("GetUserByEmail", mock.Anything, "new@example.com").Return(models.User{}, nil)
				userRepoMock.On("AddUser", mock.Anything, mock.Anything).Return(1, nil)
			},
			expectedError: false,
		},
		{
			name:      "User already exists",
			inputUser: models.User{Email: "existing@example.com"},
			mockSetup: func(userRepoMock *mocks.UserRepositoryInterface) {
				userRepoMock.On("GetUserByEmail", mock.Anything, "existing@example.com").Return(models.User{ID: 5, Email: "existing@example.com"}, nil)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize mocks using the new syntax
			UserRepoMock := mocks.NewUserRepositoryInterface(t)
			TaskRepoMock := mocks.NewTaskRepositoryInterface(t)
			logger := log.New(os.Stdout, "test: ", log.LstdFlags)

			// Initialize service with mocks and logger
			s := NewService(UserRepoMock, TaskRepoMock, logger)

			// Set up mock expectations
			tc.mockSetup(UserRepoMock)

			// Call the service method
			_, _, err := s.RegisterUser(context.TODO(), tc.inputUser)

			// Check for errors
			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	testCases := []struct {
		name          string
		inputUser     models.User
		mockSetup     func(userRepoMock *mocks.UserRepositoryInterface)
		expectedError bool
	}{
		{
			name:      "User exists",
			inputUser: models.User{Email: "test@test.com", Password: "password"},
			mockSetup: func(userRepoMock *mocks.UserRepositoryInterface) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				userRepoMock.On("GetUserByEmail", mock.Anything, "test@test.com").Return(models.User{ID: 1, Email: "test@test.com", Password: string(hashedPassword)}, nil)
			},
			expectedError: false,
		},
		{
			name:      "User does not exist",
			inputUser: models.User{Email: "test@test.com", Password: "password"},
			mockSetup: func(userRepoMock *mocks.UserRepositoryInterface) {
				userRepoMock.On("GetUserByEmail", mock.Anything, mock.Anything).Return(models.User{}, sql.ErrNoRows)
			},
			expectedError: true,
		},
		{
			name:      "Invalid password",
			inputUser: models.User{Email: "test@test.com", Password: "wrongPasswrod"},
			mockSetup: func(userRepoMock *mocks.UserRepositoryInterface) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				userRepoMock.On("GetUserByEmail", mock.Anything, "test@test.com").Return(models.User{ID: 1, Email: "test@test.com", Password: string(hashedPassword)}, nil)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			UserRepoMock := mocks.NewUserRepositoryInterface(t)
			TaskRepoMock := mocks.NewTaskRepositoryInterface(t)
			logger := log.New(os.Stdout, "test: ", log.LstdFlags)

			s := NewService(UserRepoMock, TaskRepoMock, logger)

			tc.mockSetup(UserRepoMock)

			_, _, err := s.LoginUser(context.TODO(), tc.inputUser)

			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestCheckUserAuthority(t *testing.T) {
	testCases := []struct {
		name          string
		userID        int
		taskID        int
		taskUserID    int
		mockSetup     func(taskRepoMock *mocks.TaskRepositoryInterface)
		expectedError bool
	}{
		{
			name:       "User authorized",
			userID:     1,
			taskID:     5,
			taskUserID: 1,
			mockSetup: func(taskRepoMock *mocks.TaskRepositoryInterface) {
				taskRepoMock.On("GetTaskByID", context.TODO(), mock.AnythingOfType("int")).Return(models.Task{UserID: 1, Title: "Test Task", ID: 5}, nil)
			},
			expectedError: false,
		},
		{
			name:       "User not authorized",
			userID:     1,
			taskID:     5,
			taskUserID: 2,
			mockSetup: func(taskRepoMock *mocks.TaskRepositoryInterface) {
				taskRepoMock.On("GetTaskByID", context.TODO(), mock.AnythingOfType("int")).Return(models.Task{UserID: 2, Title: "Test task", ID: 5}, nil)
			},
			expectedError: true,
		},
		{
			name:       "Task not found",
			userID:     1,
			taskID:     5,
			taskUserID: 0,
			mockSetup: func(taskRepoMock *mocks.TaskRepositoryInterface) {
				taskRepoMock.On("GetTaskByID", context.TODO(), mock.AnythingOfType("int")).Return(models.Task{}, sql.ErrNoRows)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			UserRepoMock := mocks.NewUserRepositoryInterface(t)
			TaskRepoMock := mocks.NewTaskRepositoryInterface(t)
			logger := log.New(os.Stdout, "test: ", log.LstdFlags)

			s := NewService(UserRepoMock, TaskRepoMock, logger)

			tc.mockSetup(TaskRepoMock)

			err := s.CheckUserAuthority(context.TODO(), tc.userID, tc.taskID)

			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestAddTask(t *testing.T) {
	testCases := []struct {
		name          string
		inputTask     models.Task
		mockSetup     func(taskRepoMock *mocks.TaskRepositoryInterface)
		expectedError bool
	}{
		{
			name:      "Add task successfully",
			inputTask: models.Task{Title: "Test Task", UserID: 1},
			mockSetup: func(taskRepoMock *mocks.TaskRepositoryInterface) {
				taskRepoMock.On("AddTask", mock.Anything, mock.Anything).Return(models.Task{Title: "Test Task", UserID: 1, ID: 1}, nil)
			},
			expectedError: false,
		},
		{
			name:      "Add task with error",
			inputTask: models.Task{Title: "Test Task", UserID: 1},
			mockSetup: func(taskRepoMock *mocks.TaskRepositoryInterface) {
				taskRepoMock.On("AddTask", mock.Anything, mock.Anything).Return(models.Task{}, sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			UserRepoMock := mocks.NewUserRepositoryInterface(t)
			TaskRepoMock := mocks.NewTaskRepositoryInterface(t)
			logger := log.New(os.Stdout, "test: ", log.LstdFlags)

			s := NewService(UserRepoMock, TaskRepoMock, logger)

			tc.mockSetup(TaskRepoMock)

			_, err := s.AddTask(context.TODO(), tc.inputTask)

			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}
func TestUpdateTask(t *testing.T) {
	testCases := []struct {
		name          string
		inputTask     models.Task
		taskID        int
		userID        int
		mockSetup     func(taskRepoMock *mocks.TaskRepositoryInterface)
		expectedError bool
	}{
		{
			name:      "Update task successfully",
			inputTask: models.Task{Title: "Updated Task", UserID: 1},
			taskID:    1,
			userID:    1,
			mockSetup: func(taskRepoMock *mocks.TaskRepositoryInterface) {
				taskRepoMock.On("UpdateTask", mock.Anything, mock.Anything, 1, 1).Return(models.Task{Title: "Updated Task", UserID: 1, ID: 1}, nil)
			},
			expectedError: false,
		},
		{
			name:      "Update task with error",
			inputTask: models.Task{Title: "Updated Task", UserID: 1},
			taskID:    1,
			userID:    1,
			mockSetup: func(taskRepoMock *mocks.TaskRepositoryInterface) {
				taskRepoMock.On("UpdateTask", mock.Anything, mock.Anything, 1, 1).Return(models.Task{}, sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			UserRepoMock := mocks.NewUserRepositoryInterface(t)
			TaskRepoMock := mocks.NewTaskRepositoryInterface(t)
			logger := log.New(os.Stdout, "test: ", log.LstdFlags)

			s := NewService(UserRepoMock, TaskRepoMock, logger)

			tc.mockSetup(TaskRepoMock)

			_, err := s.UpdateTask(context.TODO(), tc.inputTask, tc.taskID, tc.userID)

			if (err != nil) != tc.expectedError {
				t.Errorf("Expected error: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}
