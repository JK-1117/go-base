package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/JK-1117/go-htmx-base/internal/controller"
	"github.com/JK-1117/go-htmx-base/internal/database"
	logging "github.com/JK-1117/go-htmx-base/internal/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

type Router struct {
	Echo         *echo.Echo
	Router       *echo.Group
	Controller   *controller.Controller
	SessionStore *SessionStore
}

func NewRouter(db *sql.DB, q *database.Queries, redis *redis.Client) *Router {
	logger, _ := logging.GetLogger()

	e := echo.New()
	// e.Pre(middleware.HTTPSRedirect())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogRequestID: true,
		LogURI:       true,
		LogMethod:    true,
		LogStatus:    true,
		LogError:     true,
		LogHost:      true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		LogLatency:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			msg := fmt.Sprintf(
				`{"start_time": "%v", "request_id": "%v", "uri": "%v", `+
					`"method": "%v", "status": "%v", "error": "%v", `+
					`"host": "%v", "remote_ip": "%v", "user_agent": "%v", `+
					`"latency": "%s"}`,
				v.StartTime.Format(time.DateTime), v.RequestID, v.URI, v.Method, v.Status,
				v.Error, v.Host, v.RemoteIP, v.UserAgent, v.Latency,
			)
			if v.Status == http.StatusInternalServerError {
				logger.Echo.Err(msg)
			} else {
				logger.Echo.Info(msg)
			}
			return nil
		},
	}))
	// e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// AllowOrigins: []string{"https://labstack.com", "https://labstack.net"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions, http.MethodHead},
	}))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "Deny",
		ContentSecurityPolicy: "default-src 'self'",
	}))

	controller := controller.NewController(db, q, redis)
	store := NewSessionStore(q, redis)

	v1Router := e.Group("/v1")

	router := Router{
		Echo:         e,
		Router:       v1Router,
		Controller:   controller,
		SessionStore: store,
	}
	router.UseAuthRoute()
	v1Router.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-XSRF-TOKEN",
		CookieSameSite: http.SameSiteLaxMode,
	}))
	v1Router.Use(store.SessionAuth)
	router.UseAccountRoute()

	return &router
}

func (router *Router) Serve(port string) {
	router.Echo.Logger.Fatal(router.Echo.Start(":" + port))
}
