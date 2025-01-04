package logger

import (
	"log"
	"os"
)

// Set custom loggers for each log level
var (
	debugLogger = log.New(os.Stdout, "[scraper-DEBUG] ", log.LstdFlags)
	infoLogger  = log.New(os.Stdout, "[scraper-INFO] ", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "[scraper-ERROR] ", log.LstdFlags)
)

func Debug(text string) {
	debugLogger.Println(text)
}

func Info(text string) {
	infoLogger.Println(text)
}

func Error(err interface{}) {
	switch v := err.(type) {
	case string:
		errorLogger.Println(v)
	case error:
		errorLogger.Println(v)
	default:
		errorLogger.Println("Unknown error type")
	}
}
