package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logger *MLogger

func GetLogger() *MLogger {
	return logger
}

type MLogger struct {
	errorLogger *log.Logger
	infoLogger  *log.Logger
	debugLogger *log.Logger
}

func NewMLogger(
	prefix string,
	debug io.Writer,
	info io.Writer,
	errors io.Writer,
) *MLogger {
	l := &MLogger{
		debugLogger: log.New(debug, prefix+" [DEBUG] ", log.Ldate|log.Ltime),
		infoLogger:  log.New(info, prefix+" [INFO] ", log.Ldate|log.Ltime),
		errorLogger: log.New(errors, prefix+" [ERROR] ", log.Ldate|log.Ltime),
	}
	return l
}

func (m MLogger) Debug(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	m.debugLogger.Printf(s)
}

func (m MLogger) Debugln(v ...any) {
	m.debugLogger.Println(v...)
}

func (m MLogger) Info(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	m.infoLogger.Printf(s)
}

func (m MLogger) Infoln(v ...any) {
	m.infoLogger.Println(v...)
}

func (m MLogger) Error(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	m.errorLogger.Printf(s)
}

func (m MLogger) Errorln(v ...any) {
	m.errorLogger.Println(v...)
}

func Default(prefix string) *MLogger {
	l := NewMLogger(prefix, os.Stdout, os.Stdout, os.Stderr)
	logger = l
	return l
}

func InitLoggerWithFiles(prefix string) (*MLogger, error) {
	now := time.Now()
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	infoFile := fmt.Sprintf(dir+string(filepath.Separator)+"%s %s-info.log", now.Format(time.DateTime))
	infoWriter, err := os.Create(infoFile)
	if err != nil {
		return nil, err
	}
	debugFile := fmt.Sprintf(dir+string(filepath.Separator)+"%s %s-debug.log", now.Format(time.DateTime))
	debugWriter, err := os.Create(debugFile)
	if err != nil {
		return nil, err
	}
	errFile := fmt.Sprintf(dir+string(filepath.Separator)+"%s %s-error.log", now.Format(time.DateTime))
	errWriter, err := os.Create(errFile)
	if err != nil {
		return nil, err
	}
	l := NewMLogger(prefix, debugWriter, infoWriter, errWriter)
	logger = l
	return l, nil
}

func InitLoggerWithFile(prefix string) (*MLogger, error) {
	now := time.Now()
	dir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	file := fmt.Sprintf(dir+string(filepath.Separator)+"%s %s.log", prefix, strings.ReplaceAll(now.Format(time.DateTime), ":", " "))
	writer, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	l := NewMLogger(prefix, writer, writer, writer)
	logger = l
	return l, nil
}
