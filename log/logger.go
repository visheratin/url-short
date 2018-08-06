// Package log implements a simple logging mechanism for
// the URL shortening microservice.
package log

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/visheratin/url-short/config"
)

// Logger is a container for various types of logs.
// Use Trace for information messages and
// Error for reports about errors in the program.
type Logger struct {
	Trace *log.Logger
	Error *log.Logger
}

// instance is a private instance of Logger.
var instance Logger

// initLog sets the path for file with logs and
// desribes a way for logs formatting.
func initLog() {
	instance = Logger{}
	var logPath string
	config, err := config.Config()
	if err != nil {
		log.Println(err)
		logPath = "."
	} else {
		logPath = config.LogPath
	}
	logFilepath := path.Join(logPath, "log.txt")
	traceW := io.Writer(os.Stdout)
	errorW := io.Writer(os.Stdout)
	file, err := os.OpenFile(logFilepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open log file")
	} else {
		traceW = io.Writer(file)
		errorW = io.Writer(file)
	}
	instance.Trace = log.New(traceW, "[TRACE]: ", log.Ldate|log.Ltime|log.Lshortfile)
	instance.Error = log.New(errorW, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Log returns a copy of a Logger instance.
func Log() Logger {
	if instance.Error == nil || instance.Trace == nil {
		initLog()
	}
	return instance
}
