package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jk1117/go-base/internal/database"
	"github.com/jk1117/go-base/internal/logger"
	"github.com/robfig/cron/v3"
)

type CronJob struct {
	q    *database.Queries
	cron *cron.Cron
}

var cFile *os.File

func NewCron(q *database.Queries) *CronJob {
	my, _ := time.LoadLocation("Asia/Kuala_Lumpur")
	Logger, _ := logger.GetLogger()

	appname := os.Getenv("APPNAME")
	cFile, err := os.OpenFile(fmt.Sprintf("./logs/%v_schedule.txt", appname), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Logger.Cron.Err(err.Error())
	}
	cWriter := io.MultiWriter(os.Stdout, cFile)
	cLogger := log.New(cWriter, "cron: ", log.LstdFlags)

	c := cron.New(
		cron.WithLogger(
			cron.VerbosePrintfLogger(cLogger),
		),
		cron.WithLocation(my),
	)

	return &CronJob{
		q:    q,
		cron: c,
	}
}

func (job *CronJob) Start() {
	job.cron.AddFunc("@daily", job.CleanInvalidSession)

	job.cron.Start()
}

func (job *CronJob) Stop() context.Context {
	cFile.Close()
	return job.cron.Stop()
}

func (job *CronJob) CleanInvalidSession() {
	Logger, _ := logger.GetLogger()
	Logger.Cron.Info(fmt.Sprintf("STARTING CleanInvalidSession"))
	if err := job.q.DeleteExpiredSession(context.Background()); err != nil {
		Logger.Cron.Err(fmt.Sprintf("CleanInvalidSession: %v", err))
		return
	}
	Logger.Cron.Info(fmt.Sprintf("COMPLETED CleanInvalidSession"))
}

func (job *CronJob) RotateLogFiles() {
	Logger, _ := logger.GetLogger()
	Logger.Cron.Info(fmt.Sprintf("STARTING RotateLogFiles"))
	if err := Logger.RotateLogFiles(); err != nil {
		Logger.Cron.Err(fmt.Sprintf("RotateLogFiles: %v", err))
		return
	}
	Logger.Cron.Info(fmt.Sprintf("COMPLETED RotateLogFiles"))
}
