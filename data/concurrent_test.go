// SPDX-License-Identifier: MIT
// Copyright © 2025 Tom Straub <github.com/straubt1>

package data

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestMapConcurrentOrderPreserved(t *testing.T) {
	items := []int{10, 20, 30, 40}
	got, err := mapConcurrent(4, items, func(n int) (int, error) {
		time.Sleep(time.Millisecond)
		return n * 2, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []int{20, 40, 60, 80}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestMapConcurrentSerial(t *testing.T) {
	var inFlight int32
	var maxInFlight int32
	items := []int{1, 2, 3, 4}
	_, err := mapConcurrent(1, items, func(n int) (int, error) {
		cur := atomic.AddInt32(&inFlight, 1)
		for {
			old := atomic.LoadInt32(&maxInFlight)
			if cur <= old || atomic.CompareAndSwapInt32(&maxInFlight, old, cur) {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
		atomic.AddInt32(&inFlight, -1)
		return n, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if max := atomic.LoadInt32(&maxInFlight); max != 1 {
		t.Fatalf("max in-flight = %d, want 1", max)
	}
}

func TestMapConcurrentFirstErrorStopsLaterItems(t *testing.T) {
	var mu sync.Mutex
	var started []int
	items := []int{0, 1, 2, 3}
	_, err := mapConcurrent(1, items, func(n int) (int, error) {
		mu.Lock()
		started = append(started, n)
		mu.Unlock()
		if n == 1 {
			return 0, errors.New("boom")
		}
		return n, nil
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "boom" {
		t.Fatalf("err = %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(started) != 2 || started[0] != 0 || started[1] != 1 {
		t.Fatalf("started = %v, want [0 1]", started)
	}
}

func TestMapConcurrentRejectsInvalidConcurrency(t *testing.T) {
	_, err := mapConcurrent(0, []int{1}, func(n int) (int, error) { return n, nil })
	if err == nil {
		t.Fatal("expected error for concurrency 0")
	}
	_, err = mapConcurrent(-1, []int{1}, func(n int) (int, error) { return n, nil })
	if err == nil {
		t.Fatal("expected error for concurrency -1")
	}
}
