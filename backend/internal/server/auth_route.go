package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jk1117/go-base/internal/controller"
	"github.com/jk1117/go-base/internal/logger"
	"github.com/labstack/echo/v4"
)

func (r *Router) UseAuthRoute() {
	r.Router.POST("/signup", r.SignUp)

	r.Router.POST("/login", r.LogIn)
}

func (r *Router) SignUp(c echo.Context) error {
	Logger, _ := logger.GetLogger()
	decoder := json.NewDecoder(c.Request().Body)
	params := controller.SignUpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		Logger.App.Err(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
	}
	user_id, err := r.Controller.SignUp(c, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	session, err := r.SessionStore.NewSession(c, user_id)
	err = r.SessionStore.SetSessionCookie(c, session)
	if err != nil {
		// User created, but session cannot be created
		Logger.App.Err(err.Error())
		return c.String(http.StatusCreated, "Signup Successfully, you may procedd to login.")
	}

	return c.NoContent(http.StatusCreated)
}

func (r *Router) LogIn(c echo.Context) error {
	Logger, _ := logger.GetLogger()
	decoder := json.NewDecoder(c.Request().Body)
	params := controller.VerifyAccountParams{}
	err := decoder.Decode(&params)
	if err != nil {
		Logger.App.Err(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
	}
	account, err := r.Controller.VerifyAccount(c, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect email or password.")
	}
	session, err := r.SessionStore.NewSession(c, account.ID)
	err = r.SessionStore.SetSessionCookie(c, session)
	if err != nil {
		Logger.App.Err(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Something went wrong, please try again later.")
	}

	return c.JSON(http.StatusOK, account)
}
