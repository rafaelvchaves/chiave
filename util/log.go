package util

import (
	"log"
	"os"
)

type Logger struct {
	iLogger, wLogger, eLogger *log.Logger
}

func NewLogger(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	l := Logger{}
	l.iLogger = log.New(file, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.wLogger = log.New(file, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	l.eLogger = log.New(file, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	return &l, nil
}

func (l *Logger) Infof(format string, args ...any) {
	l.iLogger.Printf(format, args...)
}

func (l *Logger) Warningf(format string, args ...any) {
	l.wLogger.Printf(format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.eLogger.Printf(format, args...)
}
