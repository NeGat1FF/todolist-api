package models

import "github.com/uptrace/bun"

type Task struct {
	bun.BaseModel `bun:"tasks" swaggerignore:"true"`
	ID            int    `bun:"id,pk,autoincrement" json:"id"`
	UserID        int    `bun:"user_id,notnull" json:"-"`
	Title         string `bun:"title,notnull" json:"title"`
	Description   string `bun:"description" json:"description"`
}

type TaskKey struct{}
