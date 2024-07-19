package logging

import (
	"fmt"
	"image"
	"os"
	"time"
)

type LogLevel string

const (
	System  LogLevel = "SYSTEM"
	Info    LogLevel = "INFO"
	Warning LogLevel = "WARNING"
	Error   LogLevel = "ERROR"
)

var writer = &logWriter{}

type logWriter struct {
	file *os.File
}

func (l *logWriter) Initialize() error {
	return nil
}

func (l *logWriter) Write(level LogLevel, line string) error {
	line = fmt.Sprintf("[%s] %s", level, line)
	return nil
}

type Logger struct {
	stack   []string
	loggers []*Logger
	discord *Discord
}

func (s *Logger) Log(level LogLevel, message string) error {
	levelStr := fmt.Sprintf("[%s] %s", time.Now().Format("hh:mm:ss"), message)
	return nil
}

func (s *Logger) LogUpdate(level LogLevel, message string, id *int, screenshot *image.RGBA) (*int, error) {
	return nil, nil
}

func (s *Logger) NewLogger(name string) *Logger {
	logger := &Logger{stack: append(s.stack, name), discord: s.discord}
	s.loggers = append(s.loggers, logger)
	return logger
}

func NewLogger(name string, discord *Discord) *Logger {
	return &Logger{discord: discord}
}
