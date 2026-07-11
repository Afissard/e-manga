package config

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	// Log levels (the higher the number, the more verbose)
	LogLevelError = iota
	LogLevelWarning
	LogLevelInfo

	LogFolder = "logs"
)

type Logger struct {
	LogChan      chan string
	MaxLevel     int
	logger       *log.Logger
	logFile      *os.File
	logToConsole bool
}

func (l *Logger) LogMessage(message string, level int) {
	if l == nil {
        return // If the logger is not initialized, do nothing
    }

	message = "[" + levelToString(level) + "]\t" + message
	if level <= l.MaxLevel {
		l.LogChan <- message
		if l.logToConsole {
			log.Println(message)
		}
	}
}

func InitLogger(maxLevel int) (*Logger, error) {
	// create a new log folder if it doesn't exist
	if _, err := os.Stat(LogFolder); os.IsNotExist(err) {
		err := os.Mkdir(LogFolder, 0755)
		if err != nil {
			return nil, err
		}
	}

	// create a new log file with the current date and time
	logFilePath := filepath.Join(LogFolder, time.Now().Format("2006-01-02_15-04-05")+".log")

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	logger := log.New(logFile, "", log.LstdFlags)

	return &Logger{
		LogChan:  make(chan string, 100),
		MaxLevel: maxLevel,
		logger:   logger,
		logFile:  logFile,
	}, nil
}

func (l *Logger) Close() {
	close(l.LogChan)
	l.logFile.Close()
}

func (l *Logger) Run() {
	for msg := range l.LogChan {
		l.logger.Println(msg)
	}
}

func levelToString(level int) string {
	switch level {
	case LogLevelError:
		return "ERROR"
	case LogLevelWarning:
		return "WARN"
	default:
		return "INFO"
	}
}
