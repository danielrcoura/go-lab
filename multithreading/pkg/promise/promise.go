// Package promise provides a implementation of a JavaScript-like promise in Go.
package promise

import (
	"sync"
)

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
