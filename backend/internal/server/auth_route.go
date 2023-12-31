package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JK-1117/go-base/internal/controller"
	"github.com/JK-1117/go-base/internal/database"
	"github.com/JK-1117/go-base/internal/helper"
	logging "github.com/JK-1117/go-base/internal/logger"
	"github.com/labstack/echo/v4"
)

func (r *Router) UseAuthRoute() {
	r.Router.POST("/signup", r.SignUp)
	r.Router.POST("/login", r.LogIn)
	r.Router.POST("/logout", r.LogOut)
}

func (r *Router) SignUp(c echo.Context) error {
	logger, _ := logging.GetLogger()
	decoder := json.NewDecoder(c.Request().Body)
	params := controller.SignUpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		logger.App.Err(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error parsing JSON: %v", err))
	}
	user_id, err := r.Controller.SignUp(c, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	session, err := r.SessionStore.NewSession(c, user_id)
	err = r.SessionStore.SetSessionCookie(c, session)
	if err != nil {
		// User created, but session cannot be created
		logger.App.Err(err.Error())
		return c.String(http.StatusCreated, "Signup Successfully, you may procedd to login.")
	}

	return c.NoContent(http.StatusCreated)
}

func (r *Router) LogIn(c echo.Context) error {
	logger, _ := logging.GetLogger()
	decoder := json.NewDecoder(c.Request().Body)
	params := controller.VerifyAccountParams{}
	err := decoder.Decode(&params)
	if err != nil {
		logger.App.Err(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error parsing JSON: %v", err))
	}
	account, err := r.Controller.VerifyAccount(c, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect email or password.")
	}
	session, err := r.SessionStore.NewSession(c, account.ID)
	err = r.SessionStore.SetSessionCookie(c, session)
	if err != nil {
		logger.App.Err(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrGeneralMsg)
	}

	return c.JSON(http.StatusOK, account)
}

func (r *Router) LogOut(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     SESSIONCOOKIE,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		Value:    "",
	})

	return c.NoContent(http.StatusOK)
}

func (r *Router) ForgotPassword(c echo.Context) error {
	logger, _ := logging.GetLogger()
	decoder := json.NewDecoder(c.Request().Body)
	params := controller.ForgotPasswordParams{}
	err := decoder.Decode(&params)
	if err != nil {
		logger.App.Err(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error parsing JSON: %v", err))
	}
	if err = r.Controller.ForgotPassword(c, params); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.String(http.StatusOK, "An email will be sent to your email if an account is registered under it.")
}

func (r *Router) Authorization(resource string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if resource == "" {
				return echo.NewHTTPError(http.StatusInternalServerError, "Resource missing for authorization.")
			}

			perm, err := r.Controller.GetResourcePermissions(c, controller.GetResourcePermissionsParams{
				Resource: resource,
				Roles:    c.Get(helper.C_USERROLES).([]database.RoleEnum),
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrGeneralMsg)
			}

			if perm.Read == controller.RESTRICTED {
				return echo.NewHTTPError(http.StatusForbidden, "You are not authorized to access this resource.")
			}
			c.Set(helper.C_PERMISSION, perm)
			return next(c)
		}
	}
}
