package router

import (
	"database/sql"
	"net/http"

	"github.com/JK-1117/go-base/internal/database"
	"github.com/JK-1117/go-base/internal/helper"
	"github.com/JK-1117/go-base/internal/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AccountRouter struct {
	s *services.AccountService
}

func NewAccountRouter(db *sql.DB, q *database.Queries) *AccountRouter {
	s := services.NewAccountService(db, q)

	return &AccountRouter{s: s}
}

func (r *AccountRouter) RegisterRoute(router *echo.Group, a *AuthRouter) {
	router.Use(a.Authorization("account"))
	router.GET("/me", r.Me)
}

func (r *AccountRouter) Me(c echo.Context) error {
	userId := c.Get(helper.C_USERID).(uuid.UUID)
	account, err := r.s.GetAccount(c, userId.String())
	if err != nil {
		switch v := err.(type) {
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrGeneralMsg)
		case services.UnauthorizedError:
			return echo.NewHTTPError(http.StatusUnauthorized, v.Error())
		}
	}

	return c.JSON(http.StatusOK, account)
}
