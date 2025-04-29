package parallel

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type ParamsSupplier[T any] <-chan T

type Line[T any] struct {
	done   chan struct{}
	size   uint64
	anyErr error
}

func NewLine[T any](size uint64) Line[T] {
	return Line[T]{
		done: make(chan struct{}, size),
		size: size,
	}
}

// Run continue to run until any of the following conditions are met.
// - context canceled
// - supplier closed
// - action return error or panic
func (a *Line[T]) Run(ctx context.Context, supplier ParamsSupplier[T], action func(T) error) {
	for i := uint64(0); i < a.size; i++ {
		go func() {
			defer func() {
				if msg := recover(); msg != nil {
					a.anyErr = fmt.Errorf("%v", msg)
				}
				a.done <- struct{}{}
			}()

			for params := range supplier {
				if a.isCanceled(ctx) {
					a.anyErr = errors.New("context canceled")
				}

				if a.anyErr != nil {
					a.drop(supplier)
					return
				}

				err := action(params)

				if err != nil {
					a.anyErr = err
				}
			}
		}()
	}
}

// Wait until finished
func (a *Line[T]) Wait() error {
	<-a.done
	return a.anyErr
}

// WaitTime waits for the process to finish within the given duration.
// If the process does not finish within the duration, it returns a timeout error.
func (a *Line[T]) WaitTime(timeout time.Duration) error {
	select {
	case <-a.done:
		return a.anyErr
	case <-time.After(timeout):
		return errors.New("wait timeout")
	}
}

// Error return anyErr's value
func (a *Line[T]) Error() string {
	return a.anyErr.Error()
}

func (a *Line[T]) isCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// drop clean the paramsSupplier
func (a *Line[T]) drop(supplier ParamsSupplier[T]) {
	for range supplier {
	}
}
