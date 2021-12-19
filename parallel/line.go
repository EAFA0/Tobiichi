package parallel

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type ParamsSupplier <-chan interface{}

type Line struct {
	wg     sync.WaitGroup
	size   uint64
	anyErr error
}

func NewLine(size uint64) Line {
	return Line{size: size}
}

// Run continue to run until any of the following conditions are met.
// - context canceled
// - supplier closed
// - action return error or panic
func (a *Line) Run(ctx context.Context, supplier ParamsSupplier, action func(interface{}) error) {
	for i := uint64(0); i < a.size; i++ {
		a.wg.Add(1)
		go func() {
			defer func() {
				if msg := recover(); msg != nil {
					a.anyErr = fmt.Errorf("%v", msg)
				}
				a.wg.Done()
			}()

			for params := range supplier {
				if a.isCanceld(ctx) {
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
func (a *Line) Wait() error {
	a.wg.Wait()
	return a.anyErr
}

// Error return anyErr's value
func (a *Line) Error() string {
	return a.anyErr.Error()
}

func (a *Line) isCanceld(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// drop clean the paramsSupplier
func (a *Line) drop(supplier ParamsSupplier) {
	for range supplier {
	}
}
