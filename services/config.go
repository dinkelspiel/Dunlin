package services

import "os"

type Config struct {
	MariaDBDatabaseUrl string
	RedisDatabaseUrl   string
	StorageUrl         string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		MariaDBDatabaseUrl: os.Getenv("MARIADB_DATABASE_URL"),
		RedisDatabaseUrl:   os.Getenv("REDIS_DATABASE_URL"),
		StorageUrl:         os.Getenv("STORAGE_URL"),
	}
	return cfg, nil
}
