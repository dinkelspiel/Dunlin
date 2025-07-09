package dao

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dinkelspiel/cdn/models"
)

func ScanUserSessionRow(rows *sql.Rows) (*models.UserSession, error) {
	var user_session models.UserSession
	var createdAt string
	var updatedAt sql.NullString

	if err := rows.Scan(&user_session.Id, &user_session.UserId, &user_session.Token, &updatedAt, &createdAt); err != nil {
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
			user_session.UpdatedAt = &updatedAtTime
		} else {
			user_session.UpdatedAt = nil
		}
		user_session.CreatedAt = &createdAtTime
		return &user_session, nil
	}
}

func GetUserSessionByToken(db *sql.DB, token string) (*models.UserSession, error) {
	rows, err := db.Query("SELECT id, user_id, token, updated_at, created_at FROM user_sessions WHERE token = ?", token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		user_session, err := ScanUserSessionRow(rows)
		if err != nil {
			return nil, err
		}
		user, err := GetUserById(db, user_session.UserId)
		if err != nil {
			return nil, err
		}
		user_session.User = user
		return user_session, nil
	}
	return nil, errors.New("No session was found")
}

func GetUserSessionById(db *sql.DB, id int64) (*models.UserSession, error) {
	rows, err := db.Query("SELECT id, user_id, token, updated_at, created_at FROM user_sessions WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		user_session, err := ScanUserSessionRow(rows)
		if err != nil {
			return nil, err
		}
		user, err := GetUserById(db, user_session.UserId)
		if err != nil {
			return nil, err
		}
		user_session.User = user
		return user_session, nil
	}
	return nil, errors.New("No session was found")
}

func CreateUserSession(db *sql.DB, user models.User) (*models.UserSession, error) {
	insertUser := "INSERT INTO user_sessions(user_id, token) VALUES(?, UUID())"

	res, err := db.Exec(insertUser, user.Id)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	return GetUserSessionById(db, id)
}
