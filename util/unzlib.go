package util

import (
	"github.com/NyaaPantsu/nyaa/util/log"

	"bytes"
	"compress/zlib"
	"io/ioutil"
)

// UnZlib : Is it deprecated?
func UnZlib(description []byte) (string, error) {
	if len(description) > 0 {
		b := bytes.NewReader(description)
		z, err := zlib.NewReader(b)
		if !log.CheckError(err) {
			return "", err
		}
		defer z.Close()
		p, err := ioutil.ReadAll(z)
		if !log.CheckError(err) {
			return "", err
		}
		return string(p), nil
	}
	return "", nil
}
