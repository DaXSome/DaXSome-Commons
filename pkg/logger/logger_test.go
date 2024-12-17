package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	log := NewLogger()

	assert.NotNil(t, log)

	_, err := os.Stat(log.logDir)
	assert.NoError(t, err, "Log directory should be created")
}

func TestCreateFileHandler(t *testing.T) {
	log := NewLogger()
	logType := "test"

	err := log.CreateFileHandler(logType)
	assert.NoError(t, err, "Creating file handler should not return an error")

	logPath := filepath.Join(log.logDir, logType+".log")
	_, err = os.Stat(logPath)
	assert.NoError(t, err, "Log file should be created")
}

func TestLog(t *testing.T) {
	log := NewLogger()
	logType := Default
	scope := "TestScope"
	message := "Test log message"

	log.Log(logType, scope, message)

	logPath := filepath.Join(log.logDir, logType+".log")
	_, err := os.Stat(logPath)
	assert.NoError(t, err, "Log file should exist after logging")

	data, err := os.ReadFile(logPath)
	assert.NoError(t, err, "Should be able to read the log file")
	assert.Contains(t, string(data), message, "Log file should contain the logged message")
}

func TestCreateDuplicateFileHandler(t *testing.T) {
	log := NewLogger()
	logType := "duplicateTest"

	err := log.CreateFileHandler(logType)
	assert.NoError(t, err, "First file handler creation should not return an error")

	err = log.CreateFileHandler(logType)
	assert.Error(t, err, "Creating a duplicate file handler should return an error")
}
