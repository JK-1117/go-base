package controller

import (
	"database/sql"

	"github.com/JK-1117/go-htmx-base/internal/database"
	"github.com/redis/go-redis/v9"
)

type Controller struct {
	db  *sql.DB
	q   *database.Queries
	rdb *redis.Client
}

func NewController(db *sql.DB, q *database.Queries, rdb *redis.Client) *Controller {
	return &Controller{
		db:  db,
		q:   q,
		rdb: rdb,
	}
}

type UnauthorizedError struct {
	s string
}

func (e UnauthorizedError) Error() string {
	if e.s == "" {
		return "You are not authorized to access this resource."
	}
	return e.s
}

type ValidationError struct {
	s string
}

func (e ValidationError) Error() string {
	return e.s
}
