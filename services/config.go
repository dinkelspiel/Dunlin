package services

import "os"

type Config struct {
	MariaDBDatabaseUrl string
	RedisDatabaseUrl   string
	StorageUrl         string
	AppUrl             string
	HostRoot           string

	GmailDisplayName string
	GmailEmail       string
	GmailPassword    string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		MariaDBDatabaseUrl: os.Getenv("MARIADB_DATABASE_URL"),
		RedisDatabaseUrl:   os.Getenv("REDIS_DATABASE_URL"),
		StorageUrl:         os.Getenv("STORAGE_URL"),
		AppUrl:             os.Getenv("APP_URL"),
		HostRoot:           os.Getenv("HOST_ROOT"),

		GmailDisplayName: os.Getenv("GMAIL_DISPLAY_NAME"),
		GmailEmail:       os.Getenv("GMAIL_EMAIL"),
		GmailPassword:    os.Getenv("GMAIL_PASSWORD"),
	}
	return cfg, nil
}
