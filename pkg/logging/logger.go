package logging

import (
	"fmt"
	"github.com/nosyliam/revolution/pkg/config"
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
	stack     []string
	verbosity int
	settings  *config.Settings
}

func (s *Logger) V(verbosity int) *Logger {
	return &Logger{stack: s.stack, verbosity: verbosity, settings: s.settings}
}

func (s *Logger) Child(name string) *Logger {
	return &Logger{stack: append(s.stack, name), settings: s.settings}
}

func (s *Logger) Log(level LogLevel, message string) error {
	_ = fmt.Sprintf("[%s] %s", time.Now().Format("hh:mm:ss"), message)
	return nil
}

func (s *Logger) LogDiscord(level LogLevel, message string) error {
	return s.Log(level, message)
}

func (s *Logger) LogUpdate(level LogLevel, message string, id *int, screenshot *image.RGBA) (*int, error) {
	return nil, nil
}

func NewLogger(name string, settings *config.Settings) *Logger {
	return &Logger{stack: []string{name}, settings: settings}
}