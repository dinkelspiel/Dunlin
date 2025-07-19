package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dinkelspiel/cdn/dao"
	"github.com/dinkelspiel/cdn/db"
	"github.com/dinkelspiel/cdn/models"
	"github.com/dinkelspiel/cdn/storage"
	"github.com/gin-gonic/gin"
	"golang.org/x/sys/unix"
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

func EnsureFoldersExists() error {
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

func CreateTeamProjectFolder(teamProject models.TeamProject, relativePath string) (string, error) {
	fullPath, err := GetFilePathToFileInTeamProject(teamProject, relativePath)
	if err != nil {
		return "", err
	}

	return fullPath, os.MkdirAll(fullPath, os.ModePerm)
}

type DiskStats struct {
	Total           uint64
	Free            uint64
	HostUsed        uint64
	DunlinFilesUsed uint64
	DunlinCacheUsed uint64
}

func GetHostDiskStats() DiskStats {
	config, _ := LoadConfig()

	var stat unix.Statfs_t
	err := unix.Statfs(config.HostRoot, &stat)
	if err != nil {
		panic(err)
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	hostUsed := total - free

	dunlinFilesUsed := getDirectorySize(filepath.Join(config.StorageUrl, "public"))
	dunlinCacheUsed := getDirectorySize(filepath.Join(config.StorageUrl, ".cache"))

	return DiskStats{
		Total:           total,
		Free:            free,
		HostUsed:        hostUsed - dunlinFilesUsed - dunlinCacheUsed,
		DunlinFilesUsed: dunlinFilesUsed,
		DunlinCacheUsed: dunlinCacheUsed,
	}
}

func getDirectorySize(path string) uint64 {
	var size uint64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			size += uint64(info.Size())
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return size
}

type TeamProjectSize struct {
	TeamProject models.TeamProject
	Size        uint64
}

func GetSizeOfTeamProjectsOnDisk(db *db.DB) (*[]TeamProjectSize, error) {
	teamProjects, err := dao.GetTeamProjects(db)
	if err != nil {
		return nil, err
	}

	var teamProjectsSize []TeamProjectSize

	for _, teamProject := range *teamProjects {
		path, err := GetFilePathToFileInTeamProject(teamProject, "/")
		if err != nil {
			return nil, err
		}

		size := getDirectorySize(path)
		teamProjectsSize = append(teamProjectsSize, TeamProjectSize{
			Size:        size,
			TeamProject: teamProject,
		})

	}

	return &teamProjectsSize, nil
}

func SerializeTeamProjectSize(teamProjectSize TeamProjectSize) gin.H {
	return gin.H{
		"teamProject": models.SerializeTeamProject(teamProjectSize.TeamProject),
		"size":        teamProjectSize.Size,
	}
}

func SerializeDiskStats(diskStats DiskStats) gin.H {
	return gin.H{
		"total":           diskStats.Total,
		"free":            diskStats.Free,
		"hostUsed":        diskStats.HostUsed,
		"dunlinFilesUsed": diskStats.DunlinFilesUsed,
		"dunlinCacheUsed": diskStats.DunlinCacheUsed,
	}
}
