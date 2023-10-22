package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jk1117/go-base/internal/database"
	"github.com/jk1117/go-base/internal/logger"
	"github.com/jk1117/go-base/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env != "" {
		env = "." + env
	}
	err := godotenv.Load(".env" + env)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env variable not found.")
	}
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
	client := redis.NewClient(opt)

	_, err = logger.GetLogger()
	if err != nil {
		log.Fatal(err)
	}

	cron := server.NewCron(q)
	cron.Start()
	defer cron.Stop()

	router := server.NewRouter(db, q, client)
	router.Serve(port)
}
