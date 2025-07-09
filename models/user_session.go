package models

import "time"

type UserSession struct {
	Id        *int64
	UserId    int64
	User      *User
	Token     string
	UpdatedAt *time.Time
	CreatedAt *time.Time
}
