package wait

import (
	"context"
	"sync"
	"time"

	"k8s.io/utils/clock"
)

// DelayFunc returns the next time interval to wait.
type DelayFunc func() time.Duration

// Timer takes an arbitrary delay function and returns a timer that can handle arbitrary interval changes.
// Use Backoff{...}.Timer() for simple delays and more efficient timers.
func (fn DelayFunc) Timer(c clock.Clock) Timer {
	return &variableTimer{fn: fn, new: c.NewTimer}
}

// Until takes an arbitrary delay function and runs until cancelled or the condition indicates exit. This
// offers all of the functionality of the methods in this package.
func (fn DelayFunc) Until(ctx context.Context, immediate, sliding bool, condition ConditionWithContextFunc) error {
	return loopConditionUntilContext(ctx, &variableTimer{fn: fn, new: internalClock.NewTimer}, immediate, sliding, condition)
}

// Concurrent returns a version of this DelayFunc that is safe for use by multiple goroutines that
// wish to share a single delay timer.
func (fn DelayFunc) Concurrent() DelayFunc {
	var lock sync.Mutex
	return func() time.Duration {
		lock.Lock()
		defer lock.Unlock()
		return fn()
	}
}
