package platform

import (
	"github.com/nosyliam/revolution/pkg/image"
	"github.com/nosyliam/revolution/pkg/window"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DisplayFrames(t *testing.T) {
	frames, err := WindowBackend.DisplayFrames()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, frames, 1)
	assert.Greater(t, frames[0].Height, 500)
}

func Test_Window(t *testing.T) {
	id, err := WindowBackend.OpenWindow(window.JoinOptions{})
	assert.NoError(t, err)

	err = WindowBackend.ActivateWindow(id)
	assert.NoError(t, err)

	err = WindowBackend.SetFrame(id, image.Frame{100, 100, 0, 0})
	assert.NoError(t, err)

	/*
		proc, err := os.FindProcess(id)
		if err == nil {
			_ = proc.Kill()
		}*/
}
