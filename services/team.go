package services

import (
	"errors"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/db"
	"github.com/dinkelspiel/cdn/models"
)

func CreateTeam(db *db.DB, _team models.Team) (*models.Team, error) {
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

func GetTeamBySlug(db *db.DB, teamSlug string) (*models.Team, error) {
	team, err := dao.GetTeamBySlug(db, teamSlug)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, errors.New("no team found with slug")
	}

	return team, nil
}

func GetTeamAndProjectBySlug(db *db.DB, teamSlug string, teamProjectSlug string) (*models.Team, *models.TeamProject, error) {
	team, err := dao.GetTeamBySlug(db, teamSlug)
	if err != nil {
		return nil, nil, err
	}
	if team == nil {
		return nil, nil, errors.New("no team found with slug")
	}

	teamProject, err := dao.GetTeamProjectInTeamBySlug(db, *team, teamProjectSlug)
	if err != nil {
		return nil, nil, err

	}
	if teamProject == nil {
		return nil, nil, errors.New("no team project found with slug in team")
	}

	return team, teamProject, nil
}
