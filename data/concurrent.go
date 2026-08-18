// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
)

// mapConcurrent runs fn for each item with at most concurrency workers.
// Results are returned in input order. After the first error, no new work is
// started; in-flight calls are allowed to finish.
func mapConcurrent[T, R any](concurrency int, items []T, fn func(T) (R, error)) ([]R, error) {
	if concurrency < 1 {
		return nil, errors.New("concurrency must be at least 1")
	}
	n := len(items)
	if n == 0 {
		return []R{}, nil
	}
	if concurrency > n {
		concurrency = n
	}

	results := make([]R, n)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error
	next := int64(-1)

	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				mu.Lock()
				failed := firstErr != nil
				mu.Unlock()
				if failed {
					return
				}
				i := int(atomic.AddInt64(&next, 1))
				if i >= n {
					return
				}
				r, err := fn(items[i])
				if err != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					mu.Unlock()
					return
				}
				results[i] = r
			}
		}()
	}

	wg.Wait()
	return results, firstErr
}
