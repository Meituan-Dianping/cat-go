package cat

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	logger     *log.Logger
	mu         sync.Mutex
	currentDay int
}

func createLogger() *Logger {
	now := time.Now()

	var writer = getWriterByTime(now)

	return &Logger{
		logger:     log.New(writer, "", log.LstdFlags),
		mu:         sync.Mutex{},
		currentDay: now.Day(),
	}
}

func openLoggerFile(dirPath string, time time.Time) (*os.File, error) {
	year, month, day := time.Date()
	filename := fmt.Sprintf("%s/gocat_%d_%02d_%02d.log", dirPath, year, month, day)
	return os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}

func getWriterByTime(time time.Time) io.Writer {
	var file *os.File
	var err error

	file, err = openLoggerFile(defaultLogDir, time)
	if err == nil {
		log.Printf("Logs has been redirected to the file: %s", file.Name())
		return file
	}
	log.Printf("Cannot open file in: %s", defaultLogDir)

	file, err = openLoggerFile(tmpLogDir, time)
	if err == nil {
		log.Printf("Logs has been redirected to the file: %s", file.Name())
		return file
	}
	log.Printf("Cannot open file in: %s, logger will be disabled.", tmpLogDir)

	return ioutil.Discard
}

func (l *Logger) switchLogFile(time time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.currentDay == time.Day() {
		return
	}
	l.logger.SetOutput(getWriterByTime(time))
}

func (l *Logger) write(prefix, format string, args ...interface{}) {
	now := time.Now()

	if now.Day() != l.currentDay {
		l.switchLogFile(now)
	}
	l.logger.Printf(prefix+" "+format, args...)
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.write("[Debug]", format, args...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.write("[Info]", format, args...)
}

func (l *Logger) Warning(format string, args ...interface{}) {
	l.write("[Warning]", format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.write("[Error]", format, args...)
}

var logger = createLogger()
