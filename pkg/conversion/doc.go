// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package conversion provides go object versioning.
//
// Specifically, conversion provides a way for you to define multiple versions
// of the same object. You may write functions which implement conversion logic,
// but for the fields which did not change, copying is automated. This makes it
// easy to modify the structures you use in memory without affecting the format
// you store on disk or respond to in your external API calls.
package conversion // import "openim.cc/component-base/pkg/conversion"
