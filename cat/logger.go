package cat

import (
	"log"
)

type Logger struct {

}

func (l *Logger) Debug(format string, args ...interface{}) {
	log.Printf("[Debug] " + format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	log.Printf("[Info] " + format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	log.Printf("[Warning] " + format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	log.Printf("[Error] " + format, args...)
}

var logger Logger
