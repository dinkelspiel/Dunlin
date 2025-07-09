package models

import "time"

type Team struct {
	Id        *string
	Name      string
	OwnerId   int64
	Owner     *User
	UpdatedAt *time.Time
	CreatedAt *time.Time
}
