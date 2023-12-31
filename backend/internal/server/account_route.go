package server

import (
	"net/http"

	"github.com/JK-1117/go-htmx-base/internal/controller"
	"github.com/JK-1117/go-htmx-base/internal/helper"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (r *Router) UseAccountRoute() {
	r.Router.Use(r.Authorization("account"))
	r.Router.GET("/me", r.Me)
}

func (r *Router) Me(c echo.Context) error {
	userId := c.Get(helper.C_USERID).(uuid.UUID)
	account, err := r.Controller.GetAccount(c, userId.String())
	if err != nil {
		switch v := err.(type) {
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrGeneralMsg)
		case controller.UnauthorizedError:
			return echo.NewHTTPError(http.StatusUnauthorized, v.Error())
		}
	}

	return c.JSON(http.StatusOK, account)
}
