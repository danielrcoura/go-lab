// Package promise provides a implementation of a JavaScript-like promise in Go.
package promise

import (
	"context"
	"errors"
	"sync"
)

// All returns all the results of the promises.
// If any of the promises fail, the function will imediately return.
func All(funcs ...func() (any, error)) ([]any, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(len(funcs))
	defer wg.Wait()
	results := make([]any, len(funcs))
	for i, f := range funcs {
		go func() {
			data, err := f()
			if err != nil {
				cancel()
				wg.Done()
				return
			}
			results[i] = data
			wg.Done()
		}()
	}
	wg.Wait()
	select {
	case <-ctx.Done():
		return results, errors.New("aborted")
	default:
		return results, nil
	}
}

func AllSettled(funcs ...func() string) []string {
	wg := sync.WaitGroup{}
	wg.Add(len(funcs))
	defer wg.Wait()
	results := make([]string, len(funcs))
	for i, f := range funcs {
		go func() {
			results[i] = f()
			wg.Done()
		}()
	}
	return results
}
