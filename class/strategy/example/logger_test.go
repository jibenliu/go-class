package example

import "testing"

func TestNewLogManager(t *testing.T) {
	fileLogging := &FileLogging{}
	logManager := NewLogManager(fileLogging)

	logManager.Info()
	logManager.Error()

	dbLogging := &DbLogging{}
	logManager = NewLogManager(dbLogging)

	logManager.Info()
	logManager.Error()


}
