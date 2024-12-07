package logger

import (
	"log"
	"os"

	m "own-redis/models"
)

func InitLogger() error {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	m.Logger = log.New(logFile, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}
