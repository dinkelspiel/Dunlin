package storage

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type FSType string

const (
	FSFile FSType = "file"
	FSDir  FSType = "dir"
)

type FSItem struct {
	Type         FSType    `json:"type"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
}

func ListFiles(dirPath string) ([]FSItem, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var files []FSItem
	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		info, err := os.Lstat(fullPath)
		if err != nil {
			continue
		}

		itemType := FSFile

		// Handle symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			targetInfo, err := os.Stat(fullPath)
			if err == nil && targetInfo.IsDir() {
				itemType = FSDir
			}
		} else if info.IsDir() {
			itemType = FSDir
		}

		files = append(files, FSItem{
			Type:         itemType,
			Name:         entry.Name(),
			Size:         info.Size(),
			LastModified: info.ModTime(),
		})
	}
	return files, nil
}

func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func CreateDir(filePath string) error {
	return os.Mkdir(filePath, os.ModePerm)
}

func SerializeFSItem(fsItem FSItem) gin.H {
	return gin.H{
		"type":         fsItem.Type,
		"name":         fsItem.Name,
		"size":         fsItem.Size,
		"lastModified": fsItem.LastModified.Format(time.RFC3339),
	}
}
