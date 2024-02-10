package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	mock := NewMock()
	// Must not panic
	mock.Run()
	mock.Stop()

	// Classic mock
	var called bool
	mock.RunFunc = func() { called = true }
	mock.StopFunc = func() { called = true }

	called = false
	mock.Run()
	assert.True(t, called)

	called = false
	mock.Stop()
	assert.True(t, called)
}
