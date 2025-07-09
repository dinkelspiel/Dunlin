package services

import (
	"database/sql"
	"math/rand/v2"
	"strconv"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/models"
	"github.com/gosimple/slug"
)

func RegisterUser(db *sql.DB, _user models.User) (*models.User, error) {
	// Create user
	user, err := dao.CreateUser(db, _user)
	if err != nil {
		return nil, err
	}

	// Create team for user
	name := user.Username + "'s team"
	existing_team, err := dao.GetTeamByName(db, name)
	if err != nil {
		return nil, err
	}

	if existing_team != nil {
		// TODO: make this actually do something useful like adding an integer
		// instead of just screaming like a baby when the name already exists
		name = strconv.Itoa(rand.Int())
	}

	team := models.Team{
		Name:    name,
		Slug:    slug.Make(name),
		OwnerId: *user.Id,
	}

	_, err = dao.CreateTeam(db, team)
	if err != nil {
		return nil, err
	}

	return user, nil
}
