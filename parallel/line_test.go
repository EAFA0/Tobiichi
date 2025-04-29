package parallel

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var emptyVal = struct{}{}

func TestWait(t *testing.T) {
	line, supplier := NewLine[any](2), make(chan interface{})

	line.Run(context.Background(), supplier, func(obj interface{}) error {
		return nil
	})

	for i := 0; i < 10; i++ {
		supplier <- i
	}
	close(supplier)
	assert.NoError(t, line.Wait())
}

func TestCtxCancel(t *testing.T) {
	line, supplier := NewLine[any](1), make(chan interface{})

	ctx, cancel := context.WithCancel(context.Background())
	line.Run(ctx, supplier, func(obj interface{}) error {
		return nil
	})

	cancel()
	supplier <- emptyVal
	close(supplier)

	if assert.Error(t, line.Wait()) {
		assert.Equal(t, "context canceled", line.Error())
	}
}

func TestError(t *testing.T) {
	line, supplier := NewLine[any](1), make(chan interface{})

	msg := "errmsg"
	line.Run(context.Background(), supplier, func(obj interface{}) error {
		return errors.New(msg)
	})

	supplier <- emptyVal
	close(supplier)
	if assert.Error(t, line.Wait()) {
		assert.Equal(t, msg, line.Error())
	}
}

func TestPanic(t *testing.T) {
	line, supplier := NewLine[any](1), make(chan interface{})

	msg := "panic msg"
	line.Run(context.Background(), supplier, func(obj interface{}) error {
		panic(msg)
	})

	supplier <- emptyVal
	close(supplier)
	if assert.Error(t, line.Wait()) {
		assert.Equal(t, msg, line.Error())
	}
}
