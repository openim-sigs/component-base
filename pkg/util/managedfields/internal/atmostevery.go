// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal

import (
	"sync"
	"time"
)

// AtMostEvery will never run the method more than once every specified
// duration.
type AtMostEvery struct {
	delay    time.Duration
	lastCall time.Time
	mutex    sync.Mutex
}

// NewAtMostEvery creates a new AtMostEvery, that will run the method at
// most every given duration.
func NewAtMostEvery(delay time.Duration) *AtMostEvery {
	return &AtMostEvery{
		delay: delay,
	}
}

// updateLastCall returns true if the lastCall time has been updated,
// false if it was too early.
func (s *AtMostEvery) updateLastCall() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if time.Since(s.lastCall) < s.delay {
		return false
	}
	s.lastCall = time.Now()
	return true
}

// Do will run the method if enough time has passed, and return true.
// Otherwise, it does nothing and returns false.
func (s *AtMostEvery) Do(fn func()) bool {
	if !s.updateLastCall() {
		return false
	}
	fn()
	return true
}
