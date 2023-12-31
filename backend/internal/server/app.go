package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JK-1117/go-base/internal/database"
	logging "github.com/JK-1117/go-base/internal/logger"
	"github.com/JK-1117/go-base/internal/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type App struct {
	auth    *router.AuthRouter
	account *router.AccountRouter
}

func NewApp() *App {
	db, q, rdb := initDB()

	cron := NewCron(q)
	cron.Start()
	defer cron.Stop()

	authRouter := router.NewAuthRouter(db, q, rdb)
	accountRouter := router.NewAccountRouter(db, q)

	return &App{
		auth:    authRouter,
		account: accountRouter,
	}
}

func (app *App) Run(port string) {
	e := initServer()
	v1Router := e.Group("/v1")

	app.auth.RegisterRoute(v1Router)
	// Only apply csrf and authentication after auth routes
	v1Router.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-XSRF-TOKEN",
		CookieSameSite: http.SameSiteLaxMode,
	}))
	v1Router.Use(app.auth.SessionStore.SessionAuth)
	app.account.RegisterRoute(v1Router, app.auth)

	e.Logger.Fatal(e.Start(":" + port))
}

func initServer() *echo.Echo {
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

	return e
}

func initDB() (*sql.DB, *database.Queries, *redis.Client) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}
	redisString := os.Getenv("REDIS_URL")
	if redisString == "" {
		log.Fatal("REDIS_URL is not found in the environment")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	q := database.New(db)

	opt, err := redis.ParseURL(redisString)
	if err != nil {
		log.Fatal("Could not connect to redis")
	}
	rdb := redis.NewClient(opt)

	return db, q, rdb
}
