// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package net

import (
	"testing"
)

func TestSplitSchemeNamePort(t *testing.T) {
	table := []struct {
		in                 string
		name, port, scheme string
		valid              bool
		normalized         bool
	}{
		{
			in:         "aoeu:asdf",
			name:       "aoeu",
			port:       "asdf",
			valid:      true,
			normalized: true,
		}, {
			in:         "http:aoeu:asdf",
			scheme:     "http",
			name:       "aoeu",
			port:       "asdf",
			valid:      true,
			normalized: true,
		}, {
			in:         "https:aoeu:",
			scheme:     "https",
			name:       "aoeu",
			port:       "",
			valid:      true,
			normalized: false,
		}, {
			in:         "https:aoeu:asdf",
			scheme:     "https",
			name:       "aoeu",
			port:       "asdf",
			valid:      true,
			normalized: true,
		}, {
			in:         "aoeu:",
			name:       "aoeu",
			valid:      true,
			normalized: false,
		}, {
			in:         "aoeu",
			name:       "aoeu",
			valid:      true,
			normalized: true,
		}, {
			in:    ":asdf",
			valid: false,
		}, {
			in:    "aoeu:asdf:htns",
			valid: false,
		}, {
			in:    "http::asdf",
			valid: false,
		}, {
			in:    "http::",
			valid: false,
		}, {
			in:    "",
			valid: false,
		},
	}

	for _, item := range table {
		scheme, name, port, valid := SplitSchemeNamePort(item.in)
		if e, a := item.scheme, scheme; e != a {
			t.Errorf("%q: Wanted %q, got %q", item.in, e, a)
		}
		if e, a := item.name, name; e != a {
			t.Errorf("%q: Wanted %q, got %q", item.in, e, a)
		}
		if e, a := item.port, port; e != a {
			t.Errorf("%q: Wanted %q, got %q", item.in, e, a)
		}
		if e, a := item.valid, valid; e != a {
			t.Errorf("%q: Wanted %t, got %t", item.in, e, a)
		}

		// Make sure valid items round trip through JoinSchemeNamePort
		if item.valid {
			out := JoinSchemeNamePort(scheme, name, port)
			if item.normalized && out != item.in {
				t.Errorf("%q: Wanted %s, got %s", item.in, item.in, out)
			}
			scheme, name, port, valid := SplitSchemeNamePort(out)
			if e, a := item.scheme, scheme; e != a {
				t.Errorf("%q: Wanted %q, got %q", item.in, e, a)
			}
			if e, a := item.name, name; e != a {
				t.Errorf("%q: Wanted %q, got %q", item.in, e, a)
			}
			if e, a := item.port, port; e != a {
				t.Errorf("%q: Wanted %q, got %q", item.in, e, a)
			}
			if e, a := item.valid, valid; e != a {
				t.Errorf("%q: Wanted %t, got %t", item.in, e, a)
			}
		}
	}
}
