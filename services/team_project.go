package services

import (
	"database/sql"
	"errors"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/models"
	"github.com/gosimple/slug"
)

func CreateTeamProject(db *sql.DB, teamProject models.TeamProject) (*models.TeamProject, error) {
	teamProjectSlug := slug.Make(teamProject.Name)

	existingTeamProject, err := dao.GetTeamProjectInTeamBySlug(db, *teamProject.Team, teamProjectSlug)
	if err != nil {
		return nil, err
	}
	if existingTeamProject != nil {
		return nil, errors.New("team project with slug already exist")
	}

	// TODO: User authorization

	_teamProject, err := dao.CreateTeamProject(db, teamProject)
	if err != nil {
		return nil, err
	}

	CreateTeamProjectFolder(*_teamProject)

	return _teamProject, nil
}
