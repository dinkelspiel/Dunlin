package services

import (
	"database/sql"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/models"
)

func CreateTeam(db *sql.DB, _team models.Team) (*models.Team, error) {
	// Create Team in DB
	team, err := dao.CreateTeam(db, _team)
	if err != nil {
		return nil, err
	}

	// Create team folder on disk
	err = CreateTeamFolder(*team)
	if err != nil {
		return nil, err
	}

	return team, nil
}
