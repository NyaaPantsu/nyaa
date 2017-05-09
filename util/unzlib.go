package util

import (
	"github.com/ewhal/nyaa/util/log"

	"bytes"
	"compress/zlib"
	"io/ioutil"
)

func UnZlib(description []byte) string {
	if len(description) > 0 {
		b := bytes.NewReader(description)
		z, err := zlib.NewReader(b)
		log.CheckError(err)
		defer z.Close()
		p, err := ioutil.ReadAll(z)
		log.CheckError(err)
		return string(p)
	}
	return ""
}
