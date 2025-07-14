package services

import "os"

type Config struct {
	DatabaseUrl string
	StorageUrl  string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		StorageUrl:  os.Getenv("STORAGE_URL"),
	}
	return cfg, nil
}
