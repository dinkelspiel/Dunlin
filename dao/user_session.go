package dao

import (
	"database/sql"
	"time"

	"github.com/dinkelspiel/cdn/db"
	"github.com/dinkelspiel/cdn/models"
)

func ScanUserSessionRow(rows *sql.Rows) (*models.UserSession, error) {
	var userSession models.UserSession
	var createdAt string
	var updatedAt sql.NullString

	if err := rows.Scan(&userSession.Id, &userSession.UserId, &userSession.Token, &updatedAt, &createdAt); err != nil {
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
			userSession.UpdatedAt = &updatedAtTime
		} else {
			userSession.UpdatedAt = nil
		}
		userSession.CreatedAt = &createdAtTime
		return &userSession, nil
	}
}

func GetUserSessionByToken(db *db.DB, token string) (*models.UserSession, error) {
	rows, err := db.MariaDB.Query("SELECT id, user_id, token, updated_at, created_at FROM user_sessions WHERE token = ?", token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		userSession, err := ScanUserSessionRow(rows)
		if err != nil {
			return nil, err
		}
		user, err := GetUserById(db, userSession.UserId)
		if err != nil {
			return nil, err
		}
		userSession.User = user
		return userSession, nil
	}
	return nil, nil
}

func GetUserSessionById(db *db.DB, id int64) (*models.UserSession, error) {
	rows, err := db.MariaDB.Query("SELECT id, user_id, token, updated_at, created_at FROM user_sessions WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		userSession, err := ScanUserSessionRow(rows)
		if err != nil {
			return nil, err
		}
		user, err := GetUserById(db, userSession.UserId)
		if err != nil {
			return nil, err
		}
		userSession.User = user
		return userSession, nil
	}
	return nil, nil
}

func CreateUserSession(db *db.DB, user models.User) (*models.UserSession, error) {
	insertUser := "INSERT INTO user_sessions(user_id, token) VALUES(?, UUID())"

	res, err := db.MariaDB.Exec(insertUser, user.Id)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	return GetUserSessionById(db, id)
}
