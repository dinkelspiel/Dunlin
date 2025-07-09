package dao

import (
	"database/sql"
	"errors"
	"time"

	"github.com/dinkelspiel/cdn/models"
)

func ScanUserAuthCodeRow(rows *sql.Rows) (*models.UserAuthCode, error) {
	var auth_code models.UserAuthCode
	var createdAt string
	var updatedAt sql.NullString

	if err := rows.Scan(&auth_code.Id, &auth_code.UserId, &auth_code.Code, &auth_code.Used, &updatedAt, &createdAt); err != nil {
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
			auth_code.UpdatedAt = &updatedAtTime
		} else {
			auth_code.UpdatedAt = nil
		}
		auth_code.CreatedAt = &createdAtTime
		return &auth_code, nil
	}
}

func GetUserAuthCodeByCode(db *sql.DB, code int64) (*models.UserAuthCode, error) {
	rows, err := db.Query("SELECT id, user_id, code, used, updated_at, created_at FROM user_auth_codes WHERE code = ?", code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		auth_code, err := ScanUserAuthCodeRow(rows)
		if err != nil {
			return nil, err
		}
		user, err := GetUserById(db, auth_code.UserId)
		if err != nil {
			return nil, err
		}
		auth_code.User = user
		return auth_code, nil
	}
	return nil, errors.New("No auth code was found")
}

func CreateUserAuthCode(db *sql.DB, auth_code models.UserAuthCode) (*models.UserAuthCode, error) {
	insertAuthCode := "INSERT INTO user_auth_codes(user_id, code, used) VALUES(?, ?, ?)"

	res, err := db.Exec(insertAuthCode, auth_code.UserId, auth_code.Code, auth_code.Used)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	result := auth_code
	result.Id = &id
	return &result, nil
}

func UpdateUserAuthCodeToUsed(db *sql.DB, auth_code models.UserAuthCode) error {
	updateAuthCode := "UPDATE user_auth_codes SET used = true WHERE id = ?"

	_, err := db.Exec(updateAuthCode, auth_code.Id)
	if err != nil {
		return err
	}

	return nil
}
