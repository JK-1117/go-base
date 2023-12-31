package services

import (
	"database/sql"

	"github.com/JK-1117/go-base/internal/database"
	"github.com/JK-1117/go-base/internal/helper"
	logging "github.com/JK-1117/go-base/internal/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AccountService struct {
	db *sql.DB
	q  *database.Queries
}

func NewAccountService(db *sql.DB, q *database.Queries) *AccountService {
	return &AccountService{
		db: db,
		q:  q,
	}
}

func (service *AccountService) GetAccount(c echo.Context, userId string) (*Account, error) {
	logger, _ := logging.GetLogger()
	perm := c.Get(helper.C_PERMISSION).(Permission)
	uid, err := uuid.Parse(userId)
	if err != nil {
		logger.App.Err(err.Error())
		return nil, err
	}

	if perm.Read == OWNER_ONLY && uid != c.Get(helper.C_USERID) {
		return nil, UnauthorizedError{}
	}

	account, err := service.q.GetActiveAccountById(c.Request().Context(), uid)
	if err != nil {
		logger.App.Err(err.Error())
		return nil, err
	}

	return parseAccount(account), nil
}
