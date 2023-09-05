// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package waitgroup implements SafeWaitGroup wrap of sync.WaitGroup.
// Add with positive delta when waiting will fail, to prevent sync.WaitGroup race issue.
package waitgroup // import "openim.cc/component-base/pkg/util/waitgroup"
