package repository

import (
	"context"

	"github.com/NeGat1FF/todolist-api/internal/models"
	"github.com/uptrace/bun"
)

type TaskRepositoryInterface interface {
	GetTasks(ctx context.Context, user_id int, page int, limit int) ([]models.Task, error)
	GetTaskByID(ctx context.Context, task_id int) (models.Task, error)
	AddTask(ctx context.Context, task models.Task) (models.Task, error)
	UpdateTask(ctx context.Context, task models.Task, task_id, user_id int) (models.Task, error)
	DeleteTask(ctx context.Context, task_id, user_id int) error
}

type TaskRepository struct {
	db *bun.DB
}

func NewTaskRepository(db *bun.DB) *TaskRepository {
	return &TaskRepository{db}
}

func (tr *TaskRepository) GetTasks(ctx context.Context, user_id int, page int, limit int) ([]models.Task, error) {
	var tasks []models.Task
	err := tr.db.NewSelect().Model((*models.Task)(nil)).Where("?0 = ?1", bun.Ident("user_id"), user_id).Order("id").Limit(limit).Offset((page-1)*limit).Scan(ctx, &tasks)
	return tasks, err
}

func (tr *TaskRepository) GetTaskByID(ctx context.Context, task_id int) (models.Task, error) {
	var task models.Task
	err := tr.db.NewSelect().Model(&task).Where("?0 = ?1", bun.Ident("id"), task_id).Scan(ctx, &task)
	return task, err
}

func (tr *TaskRepository) AddTask(ctx context.Context, task models.Task) (models.Task, error) {
	var retTask models.Task
	err := tr.db.NewInsert().Model(&task).Returning("*").Scan(ctx, &retTask)
	return retTask, err
}

func (tr *TaskRepository) UpdateTask(ctx context.Context, task models.Task, task_id, user_id int) (models.Task, error) {
	var retTask models.Task
	err := tr.db.NewUpdate().Model(&task).OmitZero().Where("?0 = ?1 AND ?2 = ?3", bun.Ident("user_id"), user_id, bun.Ident("id"), task_id).Returning("*").Scan(ctx, &retTask)
	return retTask, err
}

func (tr *TaskRepository) DeleteTask(ctx context.Context, task_id, user_id int) error {
	_, err := tr.db.NewDelete().Model((*models.Task)(nil)).Where("?0 = ?1 AND ?2 = ?3", bun.Ident("id"), task_id, bun.Ident("user_id"), user_id).Exec(ctx)
	return err
}
