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
	settings  config.Reactive
}

func (s *Logger) Child(name string) *Logger {
	return &Logger{stack: append(s.stack, name), settings: s.settings}
}

func (s *Logger) Log(verbosity int, level LogLevel, message string) error {
	_ = fmt.Sprintf("[%s] %s: %s", time.Now().Format("hh:mm:ss"), level, message)
	return nil
}

func (s *Logger) LogDiscord(level LogLevel, message string, id *int, screenshot *image.RGBA) (int, error) {
	return 0, nil
}

func NewLogger(name string, settings config.Reactive) *Logger {
	return &Logger{stack: []string{name}, settings: settings}
}
