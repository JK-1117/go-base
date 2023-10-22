package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *Router) UseAccountRoute() {
	r.Router.GET("/me", r.Me)
}

func (r *Router) Me(c echo.Context) error {
	// logger, _ := logger.GetLogger()
	userId := c.Get("UserId")
	isAdministrator := c.Get("IsAdministrator")
	userRoles := c.Get("UserRoles")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":          userId,
		"is_administrator": isAdministrator,
		"user_roles":       userRoles,
	})
}
