package debounce

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestDebounce_OnlyRunsOnce(t *testing.T) {
	var count atomic.Int32

	fn := New(100*time.Millisecond, func() {
		count.Add(1)
	})

	// call 5 times rapidly
	for i := 0; i < 5; i++ {
		fn()
	}

	time.Sleep(300 * time.Millisecond)

	if count.Load() != 1 {
		t.Errorf("expected 1 call, got %d", count.Load())
	}
}

func TestDebounce_ResetsTimer(t *testing.T) {
	var count atomic.Int32

	fn := New(100*time.Millisecond, func() {
		count.Add(1)
	})

	fn()
	time.Sleep(300 * time.Millisecond) // let it fire

	fn()
	time.Sleep(300 * time.Millisecond) // let it fire again

	if count.Load() != 2 {
		t.Errorf("expected 2 calls, got %d", count.Load())
	}

}
