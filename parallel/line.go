package parallel

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Line struct {
	size   uint64
	done   chan struct{}
	anyErr error
	cancel <-chan struct{}
}

func NewLine(size uint64) Line {
	return Line{
		done: make(chan struct{}, size),
		size: size,
	}
}

func NewLineWithContext(ctx context.Context, size uint64) Line {
	return Line{
		done:   make(chan struct{}, size),
		size:   size,
		cancel: ctx.Done(),
	}
}

func NewLineWithCancel(size uint64, cancel <-chan struct{}) Line {
	return Line{
		done:   make(chan struct{}, size),
		size:   size,
		cancel: cancel,
	}
}

func Run[T any](line *Line, supplier <-chan T, action func(T) error) {
	for i := uint64(0); i < line.size; i++ {
		go func() {
			defer func() {
				if msg := recover(); msg != nil {
					line.anyErr = fmt.Errorf("%v", msg)
				}
				line.done <- struct{}{}
			}()

			for params := range supplier {
				if isCanceled(line.cancel) {
					line.anyErr = errors.New("action canceled")
				}

				if line.anyErr != nil {
					dropChan(supplier)
					return
				}

				err := action(params)

				if err != nil {
					line.anyErr = err
				}
			}
		}()
	}
}

// Wait until finished
func (a *Line) Wait() error {
	<-a.done
	return a.anyErr
}

// WaitTime waits for the process to finish within the given duration.
// If the process does not finish within the duration, it returns a timeout error.
func (a *Line) WaitTime(timeout time.Duration) error {
	select {
	case <-a.done:
		return a.anyErr
	case <-time.After(timeout):
		return errors.New("wait timeout")
	}
}

// Error return anyErr's value
func (a *Line) Error() string {
	return a.anyErr.Error()
}

func isCanceled(cancel <-chan struct{}) bool {
	select {
	case <-cancel:
		return true
	default:
		return false
	}
}

// dropChan clean the chan avoid block
func dropChan[T any](supplier <-chan T) {
	for range supplier {
	}
}
