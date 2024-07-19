package logging

import "image"

var LogColors = map[LogLevel]int{}

type Discord struct {
	channel string
}

func NewDiscordWebhook(channel string) *Discord {
	return &Discord{channel: channel}
}

func (s *Logger) Log(level LogLevel, message string) error {
	return nil
}

func (s *Logger) LogUpdate(level LogLevel, message string, id *int, screenshot *image.RGBA) (*int, error) {
	return nil, nil
}
