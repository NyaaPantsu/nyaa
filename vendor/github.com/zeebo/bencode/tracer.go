package bencode

import (
	"log"
	"strings"
)

var traceDepth int

func prefix() string {
	return strings.Repeat("\t", traceDepth)
}

func un(val string) {
	traceDepth--
	log.Print(prefix(), "<-", val)
}

func trace(val string) string {
	log.Print(prefix(), "->", val)
	traceDepth++
	return val
}
