package logger

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func New(prefix string) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, prefix, log.LstdFlags|log.Lshortfile),
	}
}

func (l *Logger) Info(msg string) {
	l.Println("[INFO]", msg)
}

func (l *Logger) Error(msg string) {
	l.Println("[ERROR]", msg)
}

func (l *Logger) Debug(msg string) {
	l.Println("[DEBUG]", msg)
}

func (l *Logger) Fatal(msg string) {
	l.Println("[FATAL]", msg)
	os.Exit(1)
}
