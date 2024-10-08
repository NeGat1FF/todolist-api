package models

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"users" swaggerignore:"true"`
	ID            int    `bun:"id,pk,autoincrement" json:"id,omitempty"`
	Username      string `bun:"name,notnull" json:"username,omitempty"`
	Email         string `bun:"email,notnull,unique" json:"email"`
	Password      string `bun:"password,notnull" json:"password"`
}

type UserIDKey struct{}

type UserKey struct{}
