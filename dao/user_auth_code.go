package dao

import (
	"database/sql"
	"time"

	"github.com/dinkelspiel/cdn/db"
	"github.com/dinkelspiel/cdn/models"
)

func ScanUserAuthCodeRow(rows *sql.Rows) (*models.UserAuthCode, error) {
	var authCode models.UserAuthCode
	var createdAt string
	var updatedAt sql.NullString

	if err := rows.Scan(&authCode.Id, &authCode.UserId, &authCode.Code, &authCode.Used, &updatedAt, &createdAt); err != nil {
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
			authCode.UpdatedAt = &updatedAtTime
		} else {
			authCode.UpdatedAt = nil
		}
		authCode.CreatedAt = &createdAtTime
		return &authCode, nil
	}
}

func GetUnusedUserAuthCodeByCode(db *db.DB, code int64) (*models.UserAuthCode, error) {
	rows, err := db.MariaDB.Query("SELECT id, user_id, code, used, updated_at, created_at FROM user_auth_codes WHERE code = ? AND used = false", code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		authCode, err := ScanUserAuthCodeRow(rows)
		if err != nil {
			return nil, err
		}
		user, err := GetUserById(db, authCode.UserId)
		if err != nil {
			return nil, err
		}
		authCode.User = user
		return authCode, nil
	}
	return nil, nil
}

func CreateUserAuthCode(db *db.DB, authCode models.UserAuthCode) (*models.UserAuthCode, error) {
	insertAuthCode := "INSERT INTO user_auth_codes(user_id, code, used) VALUES(?, ?, ?)"

	res, err := db.MariaDB.Exec(insertAuthCode, authCode.UserId, authCode.Code, authCode.Used)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	result := authCode
	result.Id = &id
	return &result, nil
}

func UpdateUserAuthCodeToUsed(db *db.DB, authCode models.UserAuthCode) error {
	updateAuthCode := "UPDATE user_auth_codes SET used = true WHERE id = ?"

	_, err := db.MariaDB.Exec(updateAuthCode, authCode.Id)
	if err != nil {
		return err
	}

	return nil
}
