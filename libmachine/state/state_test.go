package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	assert.Equal(t, "", None.String())
	assert.Equal(t, "Running", Running.String())
	assert.Equal(t, "Paused", Paused.String())
	assert.Equal(t, "Saved", Saved.String())
	assert.Equal(t, "Stopped", Stopped.String())
	assert.Equal(t, "Stopping", Stopping.String())
	assert.Equal(t, "Starting", Starting.String())
	assert.Equal(t, "Error", Error.String())
	assert.Equal(t, "Timeout", Timeout.String())
}
