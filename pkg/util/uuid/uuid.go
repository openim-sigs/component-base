// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uuid

import (
	"github.com/google/uuid"

	"openim.cc/component-base/pkg/types"
)

func NewUUID() types.UID {
	return types.UID(uuid.New().String())
}
