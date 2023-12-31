package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JK-1117/go-base/internal/database"
	"github.com/JK-1117/go-base/internal/helper"
	logging "github.com/JK-1117/go-base/internal/logger"
	"github.com/JK-1117/go-base/internal/services"
	"github.com/JK-1117/go-base/internal/session"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type AuthRouter struct {
	s            *services.AuthService
	SessionStore *session.SessionStore
}

func NewAuthRouter(db *sql.DB, q *database.Queries, rdb *redis.Client) *AuthRouter {
	s := services.NewAuthService(db, q, rdb)
	store := session.NewSessionStore(q, rdb)

	return &AuthRouter{
		s:            s,
		SessionStore: store,
	}
}

func (r *AuthRouter) RegisterRoute(router *echo.Group) {
	router.POST("/signup", r.SignUp)
	router.POST("/login", r.LogIn)
	router.POST("/logout", r.LogOut)
}

func (r *AuthRouter) SignUp(c echo.Context) error {
	logger, _ := logging.GetLogger()
	decoder := json.NewDecoder(c.Request().Body)
	params := services.SignUpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		logger.App.Err(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error parsing JSON: %v", err))
	}
	user_id, err := r.s.SignUp(c, params)
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

func (r *AuthRouter) LogIn(c echo.Context) error {
	logger, _ := logging.GetLogger()
	decoder := json.NewDecoder(c.Request().Body)
	params := services.VerifyAccountParams{}
	err := decoder.Decode(&params)
	if err != nil {
		logger.App.Err(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error parsing JSON: %v", err))
	}
	account, err := r.s.VerifyAccount(c, params)
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

func (r *AuthRouter) LogOut(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     session.SESSIONCOOKIE,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
		Value:    "",
	})

	return c.NoContent(http.StatusOK)
}

func (r *AuthRouter) ForgotPassword(c echo.Context) error {
	logger, _ := logging.GetLogger()
	decoder := json.NewDecoder(c.Request().Body)
	params := services.ForgotPasswordParams{}
	err := decoder.Decode(&params)
	if err != nil {
		logger.App.Err(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("error parsing JSON: %v", err))
	}
	if err = r.s.ForgotPassword(c, params); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.String(http.StatusOK, "An email will be sent to your email if an account is registered under it.")
}

func (r *AuthRouter) Authorization(resource string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if resource == "" {
				return echo.NewHTTPError(http.StatusInternalServerError, "Resource missing for authorization.")
			}

			perm, err := r.s.GetResourcePermissions(c, services.GetResourcePermissionsParams{
				Resource: resource,
				Roles:    c.Get(helper.C_USERROLES).([]database.RoleEnum),
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrGeneralMsg)
			}

			if perm.Read == services.RESTRICTED {
				return echo.NewHTTPError(http.StatusForbidden, "You are not authorized to access this resource.")
			}
			c.Set(helper.C_PERMISSION, perm)
			return next(c)
		}
	}
}
