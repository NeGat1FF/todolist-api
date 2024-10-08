package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NeGat1FF/todolist-api/internal/models"
	"github.com/NeGat1FF/todolist-api/internal/repository"
	"github.com/NeGat1FF/todolist-api/internal/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type ServerError struct {
	Code    int
	Message string
}

func (s ServerError) Error() string {
	return fmt.Sprintf("%d: %s", s.Code, s.Message)
}

type Service struct {
	taskRep repository.TaskRepositoryInterface
	usrRep  repository.UserRepositoryInterface
	logger  *log.Logger
}

func NewService(usr repository.UserRepositoryInterface, task repository.TaskRepositoryInterface, logger *log.Logger) *Service {
	return &Service{taskRep: task, usrRep: usr, logger: logger}
}

func (s *Service) IssueAccessToken(user_id int) string {
	tokenString, err := utils.GenerateJWT(jwt.MapClaims{
		"iss":  "todolistApp",
		"uid":  user_id,
		"type": "access",
		"exp":  time.Now().Add(time.Hour * 12).Unix()})
	if err != nil {
		s.logger.Print(err)
		return ""
	}

	s.logger.Printf("Issued token succsessfully for %d", user_id)
	return tokenString
}

func (s *Service) IssueRefreshToken(user_id int) string {

	refreshTokenString, err := utils.GenerateJWT(jwt.MapClaims{
		"iss":  "todolistApp",
		"uid":  user_id,
		"type": "refresh",
		"exp":  time.Now().Add(time.Hour * 12).Unix(),
	})
	if err != nil {
		s.logger.Print(err)
		return ""
	}

	s.logger.Printf("Issued refresh token succsessfully for %d", user_id)
	return refreshTokenString
}

func (s *Service) RegisterUser(ctx context.Context, user models.User) (string, string, error) {
	usr, err := s.usrRep.GetUserByEmail(ctx, user.Email)
	if err != nil && err != sql.ErrNoRows {
		s.logger.Print(err)
		return "", "", ServerError{http.StatusInternalServerError, "internal server error"}
	}

	if usr.Email != "" {
		return "", "", ServerError{http.StatusConflict, "user with this email already exists"}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Print(err)
		return "", "", ServerError{http.StatusInternalServerError, "internal server error"}
	}

	user.Password = string(hashedPassword)

	fmt.Print(user)

	user_id, err := s.usrRep.AddUser(ctx, user)
	if err != nil {
		s.logger.Print(err)
		return "", "", ServerError{http.StatusInternalServerError, "internal server error"}
	}

	s.logger.Print("New user registered")
	return s.IssueAccessToken(user_id), s.IssueRefreshToken(user_id), nil
}

func (s *Service) LoginUser(ctx context.Context, user models.User) (string, string, error) {
	usr, err := s.usrRep.GetUserByEmail(ctx, user.Email)
	if err == sql.ErrNoRows {
		return "", "", ServerError{http.StatusUnauthorized, "Invalid email or password"}
	} else if err != nil {
		s.logger.Print(err)
		return "", "", ServerError{http.StatusInternalServerError, "internal server error"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(user.Password))
	if err != nil {
		s.logger.Print(err)
		return "", "", ServerError{http.StatusUnauthorized, "Invalid email or password"}
	}

	s.logger.Print("User loged in succsessfully")
	return s.IssueAccessToken(usr.ID), s.IssueRefreshToken(usr.ID), nil
}

func (s *Service) CheckUserAuthority(ctx context.Context, user_id, task_id int) error {
	task, err := s.taskRep.GetTaskByID(ctx, task_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ServerError{http.StatusNotFound, "task with this id not found"}
		}
		s.logger.Print(err)
		return ServerError{http.StatusInternalServerError, "internal server error"}
	}

	if task.UserID != user_id {
		return ServerError{http.StatusUnauthorized, "user is unauthorized"}
	}

	return nil
}

func (s *Service) AddTask(ctx context.Context, task models.Task) (models.Task, error) {
	task, err := s.taskRep.AddTask(ctx, task)
	if err != nil {
		s.logger.Print(err)
	}

	return task, err
}

func (s *Service) UpdateTask(ctx context.Context, task models.Task, task_id, user_id int) (models.Task, error) {
	task, err := s.taskRep.UpdateTask(ctx, task, task_id, user_id)
	if err != nil {
		s.logger.Print(err)
	}

	return task, err
}

func (s *Service) DeleteTask(ctx context.Context, task_id, user_id int) error {
	err := s.taskRep.DeleteTask(ctx, task_id, user_id)
	if err != nil {
		s.logger.Print(err)
	}
	return err
}

func (s *Service) GetTasks(ctx context.Context, user_id, page, limit int) ([]models.Task, error) {
	tasks, err := s.taskRep.GetTasks(ctx, user_id, page, limit)
	if err != nil {
		s.logger.Print(err)
	}
	return tasks, err
}
