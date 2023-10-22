package controller

import (
	"database/sql"

	"github.com/jk1117/go-base/internal/database"
)

type Controller struct {
	db *sql.DB
	q  *database.Queries
}

func NewController(db *sql.DB, q *database.Queries) *Controller {
	return &Controller{
		db: db,
		q:  q,
	}
}
