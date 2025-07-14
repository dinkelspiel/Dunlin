package storage

import (
	"os"
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
		info, err := entry.Info()
		if err != nil {
			continue // optionally log or return error
		}

		itemType := FSFile
		if info.IsDir() {
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

func SerializeFSItem(fsItem FSItem) gin.H {
	return gin.H{
		"type":         fsItem.Type,
		"name":         fsItem.Name,
		"size":         fsItem.Size,
		"lastModified": fsItem.LastModified.Format(time.RFC3339),
	}
}
