// Copyright © 2023 OpenIM-Sigs open source community. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package stringutil

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
)

func DecodeBase64(i string) ([]byte, error) {
	return ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(i)))
}
