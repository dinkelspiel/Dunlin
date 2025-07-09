package models

import "time"

type UserAuthCode struct {
	Id        *int64
	UserId    int64
	User      *User
	Code      int
	Used      bool
	UpdatedAt *time.Time
	CreatedAt *time.Time
}
