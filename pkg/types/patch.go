// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package types

// Similarly to above, these are constants to support HTTP PATCH utilized by
// both the client and server that didn't make sense for a whole package to be
// dedicated to.
type PatchType string

const (
	JSONPatchType           PatchType = "application/json-patch+json"
	MergePatchType          PatchType = "application/merge-patch+json"
	StrategicMergePatchType PatchType = "application/strategic-merge-patch+json"
	ApplyPatchType          PatchType = "application/apply-patch+yaml"
)
