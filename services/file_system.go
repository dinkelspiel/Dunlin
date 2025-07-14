package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dinkelspiel/cdn/models"
	"github.com/dinkelspiel/cdn/storage"
)

func GetTeamProjectFiles(teamProject models.TeamProject, relativePath string) ([]storage.FSItem, error) {
	config, _ := LoadConfig()
	baseDir := fmt.Sprintf("%s/public/%d/%d", config.StorageUrl, *teamProject.Team.Id, *teamProject.Id)

	fullPath := filepath.Join(baseDir, relativePath)
	fullPath = filepath.Clean(fullPath)

	if !strings.HasPrefix(fullPath, baseDir) {
		return nil, fmt.Errorf("access to path '%s' is denied", fullPath)
	}

	return storage.ListFiles(fullPath)
}

func GetFilePathToFileInTeamProject(teamProject models.TeamProject, relativePath string) (string, error) {
	config, _ := LoadConfig()
	baseDir := fmt.Sprintf("%s/public/%d/%d", config.StorageUrl, *teamProject.Team.Id, *teamProject.Id)

	fullPath := filepath.Join(baseDir, relativePath)
	fullPath = filepath.Clean(fullPath)

	if !strings.HasPrefix(fullPath, baseDir) {
		return "", fmt.Errorf("access to path '%s' is denied", fullPath)
	}

	return fullPath, nil
}

func GetTeamProjectFile(teamProject models.TeamProject, relativePath string) ([]byte, error) {
	path, err := GetFilePathToFileInTeamProject(teamProject, relativePath)
	if err != nil {
		return nil, err
	}

	return storage.ReadFile(path)
}

func EnsureFoldersExist() error {
	config, _ := LoadConfig()
	err := os.MkdirAll(config.StorageUrl, os.ModePerm)
	if err != nil {
		return err
	}

	publicDir := fmt.Sprintf("%s/public", config.StorageUrl)
	err = os.MkdirAll(publicDir, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func CreateTeamFolder(team models.Team) error {
	config, _ := LoadConfig()

	teamDir := fmt.Sprintf("%s/public/%d", config.StorageUrl, *team.Id)
	return os.MkdirAll(teamDir, os.ModePerm)
}

func CreateTeamProjectFolder(teamProject models.TeamProject) error {
	config, _ := LoadConfig()

	teamProjectDir := fmt.Sprintf("%s/public/%d/%d", config.StorageUrl, *teamProject.Team.Id, *teamProject.Id)
	return os.MkdirAll(teamProjectDir, os.ModePerm)
}
