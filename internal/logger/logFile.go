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

	m.Logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Llongfile)
	return nil
}
