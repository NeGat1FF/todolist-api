package repository

import (
	"context"

	"github.com/NeGat1FF/todolist-api/internal/models"
	"github.com/uptrace/bun"
)

type UserRepositoryInterface interface {
	AddUser(ctx context.Context, user models.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) AddUser(ctx context.Context, user models.User) (int, error) {
	var user_id int
	err := ur.db.NewInsert().Model(&user).Returning("id").Scan(ctx, &user_id)
	return user_id, err
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := ur.db.NewSelect().Model(&user).Where("?0 = ?1", bun.Ident("email"), email).Scan(ctx, &user)
	return user, err
}
