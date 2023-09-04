// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package wait

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

type errWrapper struct {
	wrapped error
}

func (w errWrapper) Unwrap() error {
	return w.wrapped
}
func (w errWrapper) Error() string {
	return fmt.Sprintf("wrapped: %v", w.wrapped)
}

type errNotWrapper struct {
	wrapped error
}

func (w errNotWrapper) Error() string {
	return fmt.Sprintf("wrapped: %v", w.wrapped)
}

func TestInterrupted(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			err:  ErrWaitTimeout,
			want: true,
		},
		{
			err:  context.Canceled,
			want: true,
		}, {
			err:  context.DeadlineExceeded,
			want: true,
		},
		{
			err:  errWrapper{ErrWaitTimeout},
			want: true,
		},
		{
			err:  errWrapper{context.Canceled},
			want: true,
		},
		{
			err:  errWrapper{context.DeadlineExceeded},
			want: true,
		},
		{
			err:  ErrorInterrupted(nil),
			want: true,
		},
		{
			err:  ErrorInterrupted(errors.New("unknown")),
			want: true,
		},
		{
			err:  ErrorInterrupted(context.Canceled),
			want: true,
		},
		{
			err:  ErrorInterrupted(ErrWaitTimeout),
			want: true,
		},

		{
			err: nil,
		},
		{
			err: errors.New("not a cancellation"),
		},
		{
			err: errNotWrapper{ErrWaitTimeout},
		},
		{
			err: errNotWrapper{context.Canceled},
		},
		{
			err: errNotWrapper{context.DeadlineExceeded},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Interrupted(tt.err); got != tt.want {
				t.Errorf("Interrupted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorInterrupted(t *testing.T) {
	internalErr := errInterrupted{}
	if ErrorInterrupted(internalErr) != internalErr {
		t.Fatalf("error should not be wrapped twice")
	}

	internalErr = errInterrupted{errInterrupted{}}
	if ErrorInterrupted(internalErr) != internalErr {
		t.Fatalf("object should be identical")
	}

	in := errors.New("test")
	actual, expected := ErrorInterrupted(in), (errInterrupted{in})
	if actual != expected {
		t.Fatalf("did not wrap error")
	}
	if !errors.Is(actual, errWaitTimeout) {
		t.Fatalf("does not obey errors.Is contract")
	}
	if actual.Error() != in.Error() {
		t.Fatalf("unexpected error output")
	}
	if !Interrupted(actual) {
		t.Fatalf("is not Interrupted")
	}
	if Interrupted(in) {
		t.Fatalf("should not be Interrupted")
	}
}
