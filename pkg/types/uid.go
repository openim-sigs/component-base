// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package types

// UID is a type that holds unique ID values, including UUIDs.  Because we
// don't ONLY use UUIDs, this is an alias to string.  Being a type captures
// intent and helps make sure that UIDs and names do not get conflated.
type UID string
