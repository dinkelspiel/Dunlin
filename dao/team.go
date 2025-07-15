package dao

import (
	"database/sql"
	"time"

	"github.com/dinkelspiel/cdn/db"
	"github.com/dinkelspiel/cdn/models"
)

func scanTeamRow(rows *sql.Rows, db *db.DB) (*models.Team, error) {
	var team models.Team
	var createdAt string
	var updatedAt sql.NullString

	if err := rows.Scan(&team.Id, &team.Name, &team.Slug, &team.OwnerId, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	createdAtTime, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return nil, err
	}
	team.CreatedAt = &createdAtTime

	if updatedAt.Valid {
		updatedAtTime, err := time.Parse("2006-01-02 15:04:05", updatedAt.String)
		if err != nil {
			return nil, err
		}
		team.UpdatedAt = &updatedAtTime
	} else {
		team.UpdatedAt = nil
	}

	team.Owner, err = GetUserById(db, team.OwnerId)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func GetTeamById(db *db.DB, id int64) (*models.Team, error) {
	rows, err := db.MariaDB.Query("SELECT id, name, slug, owner_id, created_at, updated_at FROM teams WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanTeamRow(rows, db)
	}
	return nil, nil
}

func GetTeamsByOwner(db *db.DB, owner models.User) (*[]models.Team, error) {
	rows, err := db.MariaDB.Query("SELECT id, name, slug, owner_id, created_at, updated_at FROM teams WHERE owner_id = ?", owner.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team

	for rows.Next() {
		team, err := scanTeamRow(rows, db)
		if err != nil {
			return nil, err
		}
		teams = append(teams, *team)
	}
	return &teams, nil
}

func GetTeamByName(db *db.DB, name string) (*models.Team, error) {
	rows, err := db.MariaDB.Query("SELECT id, name, slug, owner_id, created_at, updated_at FROM teams WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanTeamRow(rows, db)
	}
	return nil, nil
}

func GetTeamBySlug(db *db.DB, slug string) (*models.Team, error) {
	rows, err := db.MariaDB.Query("SELECT id, name, slug, owner_id, created_at, updated_at FROM teams WHERE slug = ?", slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanTeamRow(rows, db)
	}
	return nil, nil
}

func GetTeamMembers(db *db.DB, team models.Team) (*[]models.User, error) {
	rows, err := db.MariaDB.Query("SELECT users.id, users.username, users.email, users.updated_at, users.created_at FROM team_users LEFT JOIN users ON team_users.user_id = users.id WHERE team_users.team_id = ?", team.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user, err := ScanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return &users, nil
}

func CreateTeam(db *db.DB, team models.Team) (*models.Team, error) {
	insertTeam := "INSERT INTO teams(name, slug, owner_id) VALUES(?, ?, ?)"

	_, err := db.MariaDB.Exec(insertTeam, team.Name, team.Slug, team.OwnerId)
	if err != nil {
		return nil, err
	}

	return GetTeamByName(db, team.Name)
}
