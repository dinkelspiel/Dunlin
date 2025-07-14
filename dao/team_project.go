package dao

import (
	"database/sql"
	"time"

	"github.com/dinkelspiel/cdn/models"
)

func scanTeamProjectRow(rows *sql.Rows, db *sql.DB) (*models.TeamProject, error) {
	var teamProject models.TeamProject
	var createdAt string
	var updatedAt sql.NullString

	if err := rows.Scan(&teamProject.Id, &teamProject.Name, &teamProject.Slug, &teamProject.TeamId, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	createdAtTime, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return nil, err
	}
	teamProject.CreatedAt = &createdAtTime

	if updatedAt.Valid {
		updatedAtTime, err := time.Parse("2006-01-02 15:04:05", updatedAt.String)
		if err != nil {
			return nil, err
		}
		teamProject.UpdatedAt = &updatedAtTime
	} else {
		teamProject.UpdatedAt = nil
	}

	teamProject.Team, err = GetTeamById(db, teamProject.TeamId)
	if err != nil {
		return nil, err
	}

	return &teamProject, nil
}

func GetTeamProjectById(db *sql.DB, id string) (*models.TeamProject, error) {
	rows, err := db.Query("SELECT id, name, slug, team_id, created_at, updated_at FROM team_projects WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanTeamProjectRow(rows, db)
	}
	return nil, nil
}

func GetTeamProjectsByTeam(db *sql.DB, team models.Team) (*[]models.TeamProject, error) {
	rows, err := db.Query("SELECT id, name, slug, team_id, created_at, updated_at FROM team_projects WHERE team_id = ?", team.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teamProjects []models.TeamProject
	for rows.Next() {
		team, err := scanTeamProjectRow(rows, db)
		if err != nil {
			return nil, err
		}
		teamProjects = append(teamProjects, *team)
	}
	return &teamProjects, nil
}

func GetTeamProjectInTeamBySlug(db *sql.DB, team models.Team, slug string) (*models.TeamProject, error) {
	rows, err := db.Query("SELECT id, name, slug, team_id, created_at, updated_at FROM team_projects WHERE slug = ? AND team_id = ?", slug, team.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		return scanTeamProjectRow(rows, db)
	}
	return nil, nil
}

func CreateTeamProject(db *sql.DB, teamProject models.TeamProject) (*models.TeamProject, error) {
	insertTeamProject := "INSERT INTO team_projects(name, slug, team_id) VALUES(?, ?, ?)"

	res, err := db.Exec(insertTeamProject, teamProject.Name, teamProject.Slug, teamProject.TeamId)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	result := teamProject
	result.Id = &id
	return &result, nil
}
