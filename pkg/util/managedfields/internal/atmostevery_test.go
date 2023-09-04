// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package internal_test

import (
	"testing"
	"time"

	"github.com/openim-sigs/component-base/pkg/util/managedfields/internal"
)

func TestAtMostEvery(t *testing.T) {
	duration := time.Second
	delay := 179 * time.Millisecond
	atMostEvery := internal.NewAtMostEvery(delay)
	count := 0
	exit := time.NewTicker(duration)
	tick := time.NewTicker(2 * time.Millisecond)
	defer exit.Stop()
	defer tick.Stop()

	done := false
	for !done {
		select {
		case <-exit.C:
			done = true
		case <-tick.C:
			atMostEvery.Do(func() {
				count++
			})
		}
	}

	if expected := int(duration/delay) + 1; count > expected {
		t.Fatalf("Function called %d times, should have been called less than or equal to %d times", count, expected)
	}
}
