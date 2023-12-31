package controller

import (
	"github.com/JK-1117/go-base/internal/helper"
	logging "github.com/JK-1117/go-base/internal/logger"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (controller *Controller) GetAccount(c echo.Context, userId string) (*Account, error) {
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

	account, err := controller.q.GetActiveAccountById(c.Request().Context(), uid)
	if err != nil {
		logger.App.Err(err.Error())
		return nil, err
	}

	return parseAccount(account), nil
}
