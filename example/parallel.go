package example

import (
	"github.com/EAFA0/Tobiichi/parallel"
	"context"
	"fmt"
	"sync"
)

// just a example
type Param int
type Result int

// some task need to be execute with parallel
func task(param Param) (Result, error) {
	fmt.Printf("param is: %d\n", param)
	return Result(param), nil
}

// ParallelAction define a function for parallel actions.
func ParallelAction(ctx context.Context, params []Param) ([]Result, error) {
	results, store := make([]Result, 0, len(params)), sync.Map{}

	// wrapper task as a action
	action := func(param Param) error {
		// some multiple task
		result, err := task(param)
		store.Store(param, result)
		return err
	}

	// do parallel action and wait results
	err := wrapperParallelAction(ctx, params, action)
	store.Range(func(key, value interface{}) bool {
		results = append(results, value.(Result))
		return true
	})

	return results, err
}

// wrapperParallelAction
func wrapperParallelAction(ctx context.Context, params []Param, action func(Param) error) error {
	line, supplier := parallel.NewLine(2), make(chan interface{})

	// wrapper action and execute.
	line.Run(ctx, supplier, func(param interface{}) error {
		return action(param.(Param))
	})

	// Pass parameters.
	for _, param := range params {
		supplier <- param
	}
	close(supplier)

	return line.Wait()
}
