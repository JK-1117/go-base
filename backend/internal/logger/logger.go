package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	Cron *LoggerGroup
	Echo *LoggerGroup
	App  *LoggerGroup
}

type LoggerGroup struct {
	logger     *log.Logger
	fileLogger *log.Logger
	file       *os.File
}
type LogLevel string

const INFO LogLevel = "INFO"
const DEBUG LogLevel = "DEBUG"
const ERROR LogLevel = "ERROR"

var logger *Logger
var lock = &sync.Mutex{}

func GetLogger() (*Logger, error) {
	if logger == nil {
		lock.Lock()
		defer lock.Unlock()
		if logger == nil {
			logger = &Logger{}
			err := logger.init()
			return logger, err
		} else {
			return logger, nil
		}
	} else {
		return logger, nil
	}
}

func (l *Logger) init() error {
	err := os.MkdirAll("./logs", 0644)
	if err != nil {
		return err
	}
	appname := os.Getenv("APPNAME")

	cFile, err := os.OpenFile(fmt.Sprintf("./logs/%v_cron.txt", appname), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	cFileLogger := log.New(cFile, "cron: ", log.LstdFlags)
	cLogger := log.New(os.Stdout, "cron: ", log.LstdFlags)

	eFile, err := os.OpenFile(fmt.Sprintf("./logs/%v_echo.txt", appname), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	eFileLogger := log.New(eFile, "echo: ", log.LstdFlags)
	eLogger := log.New(os.Stdout, Blue("echo: "), log.LstdFlags)

	oFile, err := os.OpenFile(fmt.Sprintf("./logs/%v.txt", appname), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	oFileLogger := log.New(oFile, appname+": ", log.LstdFlags)
	oLogger := log.New(os.Stdout, Magenta(appname+": "), log.LstdFlags)

	l.Cron = &LoggerGroup{
		logger:     cLogger,
		fileLogger: cFileLogger,
		file:       cFile,
	}
	l.Echo = &LoggerGroup{
		logger:     eLogger,
		fileLogger: eFileLogger,
		file:       eFile,
	}
	l.App = &LoggerGroup{
		logger:     oLogger,
		fileLogger: oFileLogger,
		file:       oFile,
	}
	return nil
}

func (l *Logger) RotateLogFiles() error {
	lock.Lock()
	defer lock.Unlock()

	l.Close()

	appname := os.Getenv("APPNAME")
	if err := os.Rename(fmt.Sprintf("./logs/%v_cron.txt", appname), fmt.Sprintf("./logs/app_cron_%v.txt", time.Now())); err != nil {
		return err
	}
	if err := os.Rename(fmt.Sprintf("./logs/%v_echo.txt", appname), fmt.Sprintf("./logs/app_echo_%v.txt", time.Now())); err != nil {
		return err
	}
	if err := os.Rename(fmt.Sprintf("./logs/%v.txt", appname), fmt.Sprintf("./logs/app_%v.txt", time.Now())); err != nil {
		return err
	}

	return l.init()
}

func (l *Logger) Close() {
	l.Cron.file.Close()
	l.Echo.file.Close()
	l.App.file.Close()
}

func (l *LoggerGroup) Info(msg string) {
	l.println(INFO, msg)
}

func (l *LoggerGroup) Err(msg string) {
	l.println(ERROR, msg)
}

func (l *LoggerGroup) Debug(msg string) {
	l.println(DEBUG, msg)
}

func (l *LoggerGroup) println(level LogLevel, msg string) {
	levelStr := "[ " + string(level) + " ]"
	l.logger.Println(fmt.Sprintf("%s %s", levelStr, msg))

	if os.Getenv("APP_ENV") == "production" {
		lock.Lock()
		defer lock.Unlock()
		if level == ERROR {
			levelStr = RedBg(levelStr)
		}
		if level == DEBUG {
			levelStr = YellowBg(levelStr)
		}

		l.fileLogger.Println(fmt.Sprintf("%s %s", levelStr, msg))
	}
}
