package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

type LogType = string

const (
	Error   LogType = "error"
	Default LogType = "default"

	LogFormat = "[{{datetime}}] [{{level}}] [{{scope}}] {{message}} {{data}} {{extra}}\n"
)

type Logger struct {
	logDir   string
	handlers map[LogType]slog.Handler
	*slog.Logger
}

func NewLogger() *Logger {
	today := time.Now().Format("2006-01-02")
	logDir := filepath.Join("logs", today)

	os.MkdirAll(logDir, 0755)

	logger := &Logger{
		logDir:   logDir,
		Logger:   slog.New(),
		handlers: make(map[string]slog.Handler),
	}

	logger.CreateFileHandler(Default)
	logger.CreateFileHandler(Error)

	return logger
}

func (l *Logger) CreateFileHandler(fileName string) error {
	if _, ok := l.handlers[fileName]; ok {
		return fmt.Errorf("%s already registered!", fileName)
	}

	logPath := filepath.Join(l.logDir, fmt.Sprintf("%s.log", fileName))
	h, err := handler.NewFileHandler(logPath)
	if err != nil {
		log.Fatalf("Failed to create log file %s: %v", logPath, err)
	}

	fileFormatter := slog.NewTextFormatter()
	fileFormatter.SetTemplate(LogFormat)

	h.SetFormatter(fileFormatter)

	l.AddHandler(h)

	l.handlers[fileName] = h

	return nil
}

func (l *Logger) Log(logType LogType, scope string, stmts ...interface{}) {
	defer slog.Close()

	getLogLevel := func() slog.Level {
		switch logType {
		case Error:
			return slog.ErrorLevel

		case Default:
			return slog.DebugLevel

		default:
			return slog.InfoLevel
		}
	}

	consoleFormatter := slog.NewTextFormatter()
	consoleFormatter.SetTemplate(LogFormat)
	consoleFormatter.EnableColor = true

	consoleHandler := handler.NewConsoleHandler(slog.AllLevels)
	consoleHandler.SetFormatter(consoleFormatter)

	l.AddHandlers(consoleHandler)

	l.WithFields(slog.M{
		"scope": scope,
	}).Log(getLogLevel(), stmts...)
}
