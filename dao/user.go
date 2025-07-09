package dao

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dinkelspiel/cdn/models"
)

func ScanUserRow(rows *sql.Rows) (*models.User, error) {
	var user models.User
	var createdAt string
	var updatedAt sql.NullString

	if err := rows.Scan(&user.Id, &user.Username, &user.Email, &updatedAt, &createdAt); err != nil {
		return nil, err
	} else {
		var updatedAtTime time.Time
		if updatedAt.Valid {
			updatedAtTime, err = time.Parse("2006-01-02 15:04:05", updatedAt.String)
			if err != nil {
				return nil, err
			}
		}
		createdAtTime, err := time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil, err
		}
		if updatedAt.Valid {
			user.UpdatedAt = &updatedAtTime
		} else {
			user.UpdatedAt = nil
		}
		user.CreatedAt = &createdAtTime
		return &user, nil
	}
}

func GetUserById(db *sql.DB, userId int64) (*models.User, error) {
	rows, err := db.Query("SELECT id, username, email, updated_at, created_at FROM users WHERE id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return ScanUserRow(rows)
	}
	return nil, nil
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	rows, err := db.Query("SELECT id, username, email, updated_at, created_at FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return ScanUserRow(rows)
	}
	return nil, nil
}

func GetAmountOfUsers(db *sql.DB) (*int, error) {
	rows, err := db.Query("SELECT COUNT(id) AS count FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var count int

		if err := rows.Scan(&count); err != nil {
			return nil, err
		}

		return &count, nil
	}
	return nil, errors.New("unable to get count of users")
}

func CreateUser(db *sql.DB, user models.User) (*models.User, error) {
	insertUser := "INSERT INTO users(username, email) VALUES(?, ?)"

	res, err := db.Exec(insertUser, user.Username, user.Email)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	result := user
	result.Id = &id
	return &result, nil
}
