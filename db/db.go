package db

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type DB struct {
	MariaDB *sql.DB
	Redis   *redis.Client
}
