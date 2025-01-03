package logger

import (
	"log"
	"os"
)

// Set custom loggers for each log level
var (
	DEBUG = log.New(os.Stdout, "[scraper-DEBUG] ", log.LstdFlags)
	INFO  = log.New(os.Stdout, "[scraper-INFO] ", log.LstdFlags)
	ERROR = log.New(os.Stderr, "[scraper-ERROR] ", log.LstdFlags)
)

func Debug(text string) {
	DEBUG.Println(text)
}

func Info(text string) {
	INFO.Println(text)
}

func Error(err error) {
	ERROR.Println(err)
}
