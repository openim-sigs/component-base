// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package wait

import (
	"context"
	"errors"
)

// ErrWaitTimeout is returned when the condition was not satisfied in time.
//
// Deprecated: This type will be made private in favor of Interrupted()
// for checking errors or ErrorInterrupted(err) for returning a wrapped error.
var ErrWaitTimeout = ErrorInterrupted(errors.New("timed out waiting for the condition"))

// Interrupted returns true if the error indicates a Poll, ExponentialBackoff, or
// Until loop exited for any reason besides the condition returning true or an
// error. A loop is considered interrupted if the calling context is cancelled,
// the context reaches its deadline, or a backoff reaches its maximum allowed
// steps.
//
// Callers should use this method instead of comparing the error value directly to
// ErrWaitTimeout, as methods that cancel a context may not return that error.
//
// Instead of:
//
//	err := wait.Poll(...)
//	if err == wait.ErrWaitTimeout {
//	    log.Infof("Wait for operation exceeded")
//	} else ...
//
// Use:
//
//	err := wait.Poll(...)
//	if wait.Interrupted(err) {
//	    log.Infof("Wait for operation exceeded")
//	} else ...
func Interrupted(err error) bool {
	switch {
	case errors.Is(err, errWaitTimeout),
		errors.Is(err, context.Canceled),
		errors.Is(err, context.DeadlineExceeded):
		return true
	default:
		return false
	}
}

// errInterrupted
type errInterrupted struct {
	cause error
}

// ErrorInterrupted returns an error that indicates the wait was ended
// early for a given reason. If no cause is provided a generic error
// will be used but callers are encouraged to provide a real cause for
// clarity in debugging.
func ErrorInterrupted(cause error) error {
	switch cause.(type) {
	case errInterrupted:
		// no need to wrap twice since errInterrupted is only needed
		// once in a chain
		return cause
	default:
		return errInterrupted{cause}
	}
}

// errWaitTimeout is the private version of the previous ErrWaitTimeout
// and is private to prevent direct comparison. Use ErrorInterrupted(err)
// to get an error that will return true for Interrupted(err).
var errWaitTimeout = errInterrupted{}

func (e errInterrupted) Unwrap() error        { return e.cause }
func (e errInterrupted) Is(target error) bool { return target == errWaitTimeout }
func (e errInterrupted) Error() string {
	if e.cause == nil {
		// returns the same error message as historical behavior
		return "timed out waiting for the condition"
	}
	return e.cause.Error()
}
