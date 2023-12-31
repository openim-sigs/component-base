// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package version

import (
	"regexp"
	"strconv"
	"strings"
)

type versionType int

const (
	// Bigger the version type number, higher priority it is
	versionTypeAlpha versionType = iota
	versionTypeBeta
	versionTypeGA
)

var kubeVersionRegex = regexp.MustCompile("^v([\\d]+)(?:(alpha|beta)([\\d]+))?$")

func parseKubeVersion(v string) (majorVersion int, vType versionType, minorVersion int, ok bool) {
	var err error
	submatches := kubeVersionRegex.FindStringSubmatch(v)
	if len(submatches) != 4 {
		return 0, 0, 0, false
	}
	switch submatches[2] {
	case "alpha":
		vType = versionTypeAlpha
	case "beta":
		vType = versionTypeBeta
	case "":
		vType = versionTypeGA
	default:
		return 0, 0, 0, false
	}
	if majorVersion, err = strconv.Atoi(submatches[1]); err != nil {
		return 0, 0, 0, false
	}
	if vType != versionTypeGA {
		if minorVersion, err = strconv.Atoi(submatches[3]); err != nil {
			return 0, 0, 0, false
		}
	}
	return majorVersion, vType, minorVersion, true
}

// CompareKubeAwareVersionStrings compares two kube-like version strings.
// Kube-like version strings are starting with a v, followed by a major version, optional "alpha" or "beta" strings
// followed by a minor version (e.g. v1, v2beta1). Versions will be sorted based on GA/alpha/beta first and then major
// and minor versions. e.g. v2, v1, v1beta2, v1beta1, v1alpha1.
func CompareKubeAwareVersionStrings(v1, v2 string) int {
	if v1 == v2 {
		return 0
	}
	v1major, v1type, v1minor, ok1 := parseKubeVersion(v1)
	v2major, v2type, v2minor, ok2 := parseKubeVersion(v2)
	switch {
	case !ok1 && !ok2:
		return strings.Compare(v2, v1)
	case !ok1 && ok2:
		return -1
	case ok1 && !ok2:
		return 1
	}
	if v1type != v2type {
		return int(v1type) - int(v2type)
	}
	if v1major != v2major {
		return v1major - v2major
	}
	return v1minor - v2minor
}
